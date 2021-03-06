// Copyright 2019 The Vearch Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package client

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/spf13/cast"
	pkg "github.com/vearch/vearch/proto"
	"github.com/vearch/vearch/proto/request"
	"github.com/vearch/vearch/proto/response"
)

type multipleSpaceSender struct {
	senders []*spaceSender
}

func (this *multipleSpaceSender) MSearchIDs(req *request.SearchRequest) (result response.SearchResponses) {
	var wg sync.WaitGroup
	respChain := make(chan struct {
		reponse response.SearchResponses
		sender  *spaceSender
	}, len(this.senders))

	for _, s := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					//fmt.Println(r)
					respChain <- struct {
						reponse response.SearchResponses
						sender  *spaceSender
					}{reponse: response.SearchResponses{newSearchResponseWithError(s.db, s.space, 0, fmt.Errorf(cast.ToString(r)))}, sender: s}
				}
			}()
			respChain <- struct {
				reponse response.SearchResponses
				sender  *spaceSender
			}{reponse: s.MSearchIDs(req), sender: s}
		}(s)
	}

	wg.Wait()
	close(respChain)

	for r := range respChain {
		if result == nil {
			result = r.reponse
			continue
		}

		if err := r.sender.mergeResultArr(result, r.reponse, req); err != nil {
			return response.SearchResponses{newSearchResponseWithError(r.sender.db, r.sender.space, 0, err)}
		}
	}
	return result
}

func (this *multipleSpaceSender) MSearchForIDs(req *request.SearchRequest) (result []byte, err error) {
	var wg sync.WaitGroup
	respChain := make(chan struct {
		reponse []byte
		sender  *spaceSender
		err     error
	}, len(this.senders))

	for _, s := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					respChain <- struct {
						reponse []byte
						sender  *spaceSender
						err     error
					}{reponse: nil, sender: s, err: fmt.Errorf(cast.ToString(r))}
				}
			}()

			if bs, err := s.MSearchForIDs(req); err != nil {
				respChain <- struct {
					reponse []byte
					sender  *spaceSender
					err     error
				}{reponse: nil, sender: s, err: nil}
			} else {
				respChain <- struct {
					reponse []byte
					sender  *spaceSender
					err     error
				}{reponse: bs, sender: s, err: nil}
			}

		}(s)
	}

	wg.Wait()
	close(respChain)

	buf := bytes.Buffer{}

	for r := range respChain {

		if r.err != nil {
			return nil, r.err
		}

		if buf.Len() != 0 {
			buf.WriteString("\n")
		}
		buf.Write(r.reponse)
	}
	return buf.Bytes(), nil
}

func (this *multipleSpaceSender) MSearch(req *request.SearchRequest) (result response.SearchResponses) {
	var wg sync.WaitGroup
	respChain := make(chan struct {
		reponse response.SearchResponses
		sender  *spaceSender
	}, len(this.senders))

	for _, s := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					respChain <- struct {
						reponse response.SearchResponses
						sender  *spaceSender
					}{reponse: response.SearchResponses{newSearchResponseWithError(s.db, s.space, 0, fmt.Errorf(cast.ToString(r)))}, sender: s}
				}
			}()
			respChain <- struct {
				reponse response.SearchResponses
				sender  *spaceSender
			}{reponse: s.MSearch(req), sender: s}
		}(s)
	}

	wg.Wait()
	close(respChain)

	for r := range respChain {
		if result == nil {
			result = r.reponse
			continue
		}

		if err := r.sender.mergeResultArr(result, r.reponse, req); err != nil {
			return response.SearchResponses{newSearchResponseWithError(r.sender.db, r.sender.space, 0, err)}
		}
	}
	return result
}

func (this *multipleSpaceSender) DeleteByQuery(req *request.SearchRequest) *response.Response {
	var wg sync.WaitGroup
	respChain := make(chan *response.Response, len(this.senders))

	for _, s := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					respChain <- &response.Response{Status: pkg.ERRCODE_INTERNAL_ERROR, Err: fmt.Errorf(cast.ToString(r))}
				}
			}()
			respChain <- s.DeleteByQuery(req)
		}(s)
	}

	wg.Wait()
	close(respChain)

	var result *response.Response
	for r := range respChain {
		if r.Err != nil {
			return r
		}
		if result == nil {
			result = r
		}
	}
	return result
}

func (this *multipleSpaceSender) MSearchNew(req *request.SearchRequest) (result response.SearchResponses) {
	var wg sync.WaitGroup
	respChain := make(chan struct {
		reponse response.SearchResponses
		sender  *spaceSender
	}, len(this.senders))

	for _, s := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					respChain <- struct {
						reponse response.SearchResponses
						sender  *spaceSender
					}{reponse: response.SearchResponses{newSearchResponseWithError(s.db, s.space, 0, fmt.Errorf(cast.ToString(r)))}, sender: s}
				}
			}()
			respChain <- struct {
				reponse response.SearchResponses
				sender  *spaceSender
			}{reponse: s.MSearchNew(req), sender: s}
		}(s)
	}

	wg.Wait()
	close(respChain)

	for r := range respChain {
		if result == nil {
			result = r.reponse
			continue
		}

		if err := r.sender.mergeResultArr(result, r.reponse, req); err != nil {
			return response.SearchResponses{newSearchResponseWithError(r.sender.db, r.sender.space, 0, err)}
		}
	}
	return result
}

func (this *multipleSpaceSender) Search(req *request.SearchRequest) *response.SearchResponse {
	var wg sync.WaitGroup
	respChain := make(chan *response.SearchResponse, len(this.senders))

	for _, sender := range this.senders {
		wg.Add(1)
		go func(par *spaceSender) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
					respChain <- newSearchResponseWithError(par.db, par.space, 0, fmt.Errorf(cast.ToString(r)))
				}
			}()
			respChain <- par.Search(req)
		}(sender)
	}

	wg.Wait()
	close(respChain)

	sortOrder, err := req.SortOrder()
	if err != nil {
		return newSearchResponseWithError(this.senders[0].db, this.senders[0].space, 0, err)
	}

	var first *response.SearchResponse

	for r := range respChain {
		if first == nil {
			first = r
			continue
		}

		err := first.Merge(r, sortOrder, req.From, *req.Size)
		if err != nil {
			return newSearchResponseWithError(this.senders[0].db, this.senders[0].space, 0, err)
		}
	}

	return first
}

func (this *multipleSpaceSender) StreamSearch(req *request.SearchRequest) (dsr *response.DocStreamResult) {
	ctx := req.Context().GetContext()
	dsr = response.NewDocStreamResult(ctx)

	go func() {
		defer func() {
			dsr.AddDoc(nil)
		}()
		for _, s := range this.senders {
			spaceDsr := s.StreamSearch(req)
			for {
				doc, err := spaceDsr.Next()
				if err != nil {
					dsr.AddErr(err)
					return
				}
				if doc == nil {
					break
				}
				dsr.AddDoc(doc)
			}
		}
	}()
	return dsr
}
