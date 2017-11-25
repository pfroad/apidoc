// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
//
// 目前支持以下三种渲染方式：
//  - html: 以 HTML 格式输出文本，模板可自定义；
//  - html+: HTML 的调试模式，程序不会输出任何，而是在浏览器中展示相关页面；
//  - json: 以 JSON 格式输出内容。
package output

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/locale"
	"github.com/issue9/utils"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	Dir      string        `yaml:"dir"` // 文档的保存目录
	Elapsed  time.Duration `yaml:"-"`   // 编译用时
	ErrorLog *log.Logger   `yaml:"-"`   // 错误信息输出通道，在 html+ 模式下会用到。
}

// Init 对 Options 作一些初始化操作。
func (o *Options) Init() *app.OptionsError {
	if len(o.Dir) == 0 {
		return &app.OptionsError{Field: "dir", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if !utils.FileExists(o.Dir) {
		if err := os.MkdirAll(o.Dir, os.ModePerm); err != nil {
			msg := locale.Sprintf(locale.ErrMkdirError, err)
			return &app.OptionsError{Field: "dir", Message: msg}
		}
	}

	return nil
}

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs *doc.Doc, o *Options) error {
	// 输出之前进行一次排序，可以保证每次渲染的数据都量样的。
	sort.SliceStable(docs.Apis, func(i, j int) bool {
		return docs.Apis[i].URL < docs.Apis[j].URL
	})

	return render(docs, o)
}
