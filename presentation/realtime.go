package presentation

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

type itemsMsg []list.Item

type RealtimeSync struct {
	p      *tea.Program
	repo   data.Repository
	ctx    context.Context
	ticker *time.Ticker
}

func NewRealtimeSync(p *tea.Program, repo data.Repository, ctx context.Context, ticker *time.Ticker) RealtimeSync {
	return RealtimeSync{
		p: p, repo: repo, ctx: ctx, ticker: ticker,
	}
}

func (r *RealtimeSync) KeepSynched() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-r.ticker.C:
			actionables, err := r.repo.List()
			if err != nil {
				panic(err)
			}
			listItems := apply.ToSlice(actionables, func(a types.Actionable) list.Item {
				return UIItem{Actionable: a}
			})
			r.p.Send(itemsMsg(listItems))
		}
	}
}
