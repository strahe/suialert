//  Copyright hyperjumptech/grule-rule-engine Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package benchmark

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/strahe/suialert/types"
)

/**
  Benchmarking `engine.Execute` function by running 100 and 1000 rules with different N values
  Please refer docs/benchmarking_en.md for more info
*/

var knowledgeBase *ast.KnowledgeBase

func Benchmark_Grule_Execution_Engine(b *testing.B) {
	rules := []struct {
		name string
		fun  func()
	}{
		{"1000 rules", generateRules(1000)},
		{"10000 rules", generateRules(10000)},
		{"100000 rules", generateRules(100000)},
	}
	for _, rule := range rules {
		for k := 0; k < 10; k++ {
			b.Run(fmt.Sprintf("%s", rule.name), func(b *testing.B) {
				rule.fun()
				owner := types.HexToAddress("0xfD9A15b8a0d7EBAB0AdB0c69F1c03b05cc2a5a32")
				for i := 0; i < b.N; i++ {
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
					e := engine.NewGruleEngine()
					//Fact1
					dataCtx := ast.NewDataContext()
					err := dataCtx.Add("Event", &f1)
					if err != nil {
						b.Fail()
					}
					err = e.Execute(dataCtx, knowledgeBase)
					if err != nil {
						fmt.Print(err)
					}
				}
			})
		}
	}
}

func generateRules(count int) func() {
	return func() {
		var rb bytes.Buffer
		for i := 1; i <= count; i++ {
			_, err := rb.WriteString(MakeRule(i))
			if err != nil {
				return
			}
		}
		lib := ast.NewKnowledgeLibrary()
		ruleBuilder := builder.NewRuleBuilder(lib)
		_ = ruleBuilder.BuildRuleFromResource("exec_rules_test", "0.1.1", pkg.NewBytesResource(rb.Bytes()))
		knowledgeBase = lib.NewKnowledgeBaseInstance("exec_rules_test", "0.1.1")
	}
}
