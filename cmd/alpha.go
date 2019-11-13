package main

import (
	"log"

	"github.com/piquette/edgr/core/model"
)

func (e *Edgr) executeAlpha() {
	//
	found, filers, err := e.FilerDao.GetSet(conf.LetterStart)
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		return
	}

	hitSymbol := false
	for _, filer := range filers {
		if conf.SymbolStart != "" && filer.Symbol != conf.SymbolStart && !hitSymbol {
			continue
		}
		hitSymbol = true
		e.BackfillFiler(*filer)
	}
}

// BackfillFiler backfills a whole filer history.
func (e *Edgr) BackfillFiler(filer model.Filer) {
	e.backfillFilerWithCutoff(filer, nil)
}
