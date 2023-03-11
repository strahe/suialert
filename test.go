package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/samber/lo"
	"github.com/strahe/suialert/benchmark"
	"github.com/strahe/suialert/types"
)

func main() {

	//if getRules("A", 20) != nil {
	//	return
	//}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)

	e := engine.NewGruleEngine()
	//Fact1

	// prepare
	io := map[int]types.Address{}
	{
		start := time.Now()
		for i := 1; i <= 10000; i++ {
			owner := types.Address{}
			owner.SetBytes([]byte(lo.RandomString(40, lo.LettersCharset)))
			io[i] = owner
			err := ruleBuilder.BuildRuleFromResource(owner.Hex(), "0.1.1", getRules(owner.Hex(), rand.Intn(20)))
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("build time:", time.Since(start))
	}

	start := time.Now()
	for _, addr := range io {
		match(addr, e, lib)
	}
	fmt.Println("total took", time.Since(start))

	//err = e.Execute(dataCtx, knowledgeBase)
	//if err != nil {
	//}
}

func match(owner types.Address, e *engine.GruleEngine, lib *ast.KnowledgeLibrary) {
	knowledgeBase := lib.NewKnowledgeBaseInstance(owner.Hex(), "0.1.1")

	dataCtx := ast.NewDataContext()
	f1 := types.CoinBalanceChange{
		PackageId:         "0x0000000000000000000000000000000000000002",
		TransactionModule: "transfer_object",
		Sender:            "0x7bcb60878fb8e28d4412324842351e7261e072ec",
		ChangeType:        "Receive",
		Owner: &types.ObjectOwner{
			ObjectOwnerInternal: &types.ObjectOwnerInternal{
				AddressOwner: &owner,
			},
		},
		CoinType:     "0x2::sui::SUI",
		CoinObjectId: "0x7cf75ee1856a0ef9e6f262209420e6ea088d0edb",
		Version:      rand.Int63n(1000000),
		Amount:       rand.Int63n(10000000000),
	}
	err := dataCtx.Add("Event", &f1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = e.FetchMatchingRules(dataCtx, knowledgeBase)
	if err != nil {
		fmt.Print(err)
	}
}

func getRules(name string, count int) pkg.Resource {
	var rb bytes.Buffer
	for i := 1; i <= count; i++ {
		_, err := rb.WriteString(benchmark.MakeRule("A"+name, i))
		if err != nil {
			return nil
		}
	}
	return pkg.NewBytesResource(rb.Bytes())
}
