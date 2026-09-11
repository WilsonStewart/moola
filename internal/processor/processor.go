package processor

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

type Processor struct {
	Datafile   *Datafile
	Mu         sync.RWMutex
	sseMu      sync.Mutex
	sseClients map[chan string]struct{}
}

type Datafile struct {
	DirectiveAliases           map[string]string
	Accounts                   map[string]*Account
	AccountAliases             map[string]string
	Symbols                    SymbolsStore
	defaultEnvelopeAccountName string
}

type SymbolsStore struct {
	DirectiveNames      []string
	DirectiveAliasNames []string
	AccountNames        []string
	AccountAliasNames   []string
}

func NewProcessor() *Processor {
	processor := &Processor{
		Datafile:   NewDatafile(),
		sseClients: make(map[chan string]struct{}),
	}

	return processor
}

func NewDatafile() *Datafile {
	datafile := &Datafile{
		DirectiveAliases: make(map[string]string),
		Accounts:         make(map[string]*Account),
		AccountAliases:   make(map[string]string),
	}

	for _, directiveName := range builtinDirectiveNames {
		datafile.Symbols.DirectiveNames = append(datafile.Symbols.DirectiveNames, directiveName)
	}

	datafile.defaultEnvelopeAccountName = "to_be_assigned"

	defaultEnvelopeAccountAliasName := "tba"

	datafile.OpenAccount(OpenAccountNode{
		accountKind:           "envelope",
		requestedAccountName:  "to_be_assigned",
		requestedAccountAlias: &defaultEnvelopeAccountAliasName,
	})

	return datafile
}

func VerifyExpectedDirectiveAndFieldRange(
	rawFields []string,
	expectedDirective string,
	expectedMinimumFieldCount int,
	expectedMaximumFieldCount int,
) error {
	if rawFields[0] != expectedDirective {
		return fmt.Errorf("invalid directive %q: expected %q", rawFields[0], expectedDirective)
	}

	count := len(rawFields)
	if count < expectedMinimumFieldCount || count > expectedMaximumFieldCount {
		return fmt.Errorf("invalid field count %d for %q: expected between %d and %d",
			count, expectedDirective, expectedMinimumFieldCount, expectedMaximumFieldCount)
	}

	return nil
}

func (p *Processor) ProcessEvery1Second(ctx context.Context) {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.ReadAndProcessMasterMoo("master.moo")
			fmt.Println(p.Datafile)
		}
	}
}

func (p *Processor) ProcessMooFileChanges() {
	lastCheck := time.Now()

	for range time.Tick(300 * time.Millisecond) {
		entries, _ := os.ReadDir(".")
		triggered := false

		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".moo" {
				info, err := e.Info()
				if err == nil && info.ModTime().After(lastCheck) {
					triggered = true
					break
				}
			}
		}

		if triggered {
			if err := p.ReadAndProcessMasterMoo("master.moo"); err != nil {
				fmt.Println(err)
			}
			lastCheck = time.Now()
		}
	}
}

func (p *Processor) ReadAndProcessMasterMoo(path string) error {
	slog.Info("Processing master.moo...")

	ndf := NewDatafile()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) == 0 {
			continue
		}

		lines = append(lines, fields)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range lines {
		directive := strings.ToLower(line[0])

		if slices.Contains(ndf.Symbols.DirectiveAliasNames, directive) {
			directive = ndf.DirectiveAliases[directive]
			line[0] = directive
		}

		if !slices.Contains(ndf.Symbols.DirectiveNames, directive) {
			return fmt.Errorf("unknown directive: %q is not a known directive or directive alias", directive)
		}

		switch directive {
		case "openaccount":
			node := OpenAccountNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := ndf.OpenAccount(node); err != nil {
				return err
			}

		case "aliasdirective":
			node := AliasDirectiveNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := ndf.AliasDirective(node); err != nil {
				return err
			}

		case "assert":
			node := AssertNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := ndf.Assert(node); err != nil {
				return err
			}

		case "transact":
			node := TransactNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := ndf.Transact(node); err != nil {
				return err
			}

		case "move":
			node := MoveNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := ndf.Move(node); err != nil {
				return err
			}
		}
	}

	// datafile, err := json.MarshalIndent(ndf, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// symbols, err := json.MarshalIndent(p.Symbols, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(datafile))
	// fmt.Println(string(symbols))

	p.Mu.Lock()
	p.Datafile = ndf
	p.Mu.Unlock()

	p.refreshSSEClients()

	return nil
}

func (p *Processor) refreshSSEClients() {
	p.sseMu.Lock()
	defer p.sseMu.Unlock()
	for ch := range p.sseClients {
		select {
		case ch <- "refresh":
		default:
		}
	}
}

func (p *Processor) SSEEventsEndpointHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan string, 1)

	p.sseMu.Lock()
	p.sseClients[ch] = struct{}{}
	p.sseMu.Unlock()

	defer func() {
		p.sseMu.Lock()
		delete(p.sseClients, ch)
		p.sseMu.Unlock()
		close(ch)
	}()

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
