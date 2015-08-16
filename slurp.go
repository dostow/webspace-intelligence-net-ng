// +build slurp

package main

//This file will be only compiled along the project with slurp. So don't put any projec code here.

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/stages/archive"
	"github.com/omeid/slurp/stages/fs"
	"github.com/omeid/slurp/stages/web"

	"github.com/slurp-contrib/gin"
	"github.com/slurp-contrib/livereload"
	"github.com/slurp-contrib/resources"
	"github.com/slurp-contrib/watch"
)

func init() {
	config.Livereload = ":35729"
}

func Slurp(b *slurp.Build) {

	//Set the build meta information.
	b.Name = "IDCGroup Ltd"
	b.Usage = "Static site compiler for idcgroupltd."
	b.Version = "v0.0.1"
	b.Author = "osiloke emoekpere"
	b.Email = "osi@progwebtech.com"

	b.Task(slurp.Task{
		Name:  "libs",
		Usage: "Download frontend dependencies.",
		Action: func(c *slurp.C) error {
			return web.Get(c).Then(
				archive.Unzip(c),
				fs.Dest(c, "libs/"),
			)
		},
	},

		slurp.Task{
			Name:  "assets",
			Usage: "Concat frontend JavaScript libraries into lib.js",
			Action: func(c *slurp.C) error {
				return fs.Src(c,
					"assets/*/*/*/*",
				).Then(
					fs.Dest(c, "./public/assets/"),
				)
			},
		},

		slurp.Task{
			Name:  "index",
			Usage: "Build gcss files into style.css",
			Action: func(c *slurp.C) error {
				return fs.Src(c, "index.html").Then(
					fs.Dest(c, "./public/"),
				)
			},
		},

		slurp.Task{
			Name:  "gin",
			Usage: "Run the Gin build server and proxy.",
			Description: `Gin task uses the slurp tag to allow for package configuration.
		It sets config.Livereload to livereload port, useful for including the livereload javascript client from the template.`,
			Action: func(c *slurp.C) error {
				gin := gin.NewGin(c, &gin.Config{}, "-tags=slurp")
				watch := watch.Watch(c, gin.Run, "*.go", "*/*.go", "*/*/*.go")

				<-c.Done()
				watch.Close()
				gin.Close()
				return nil
			},
		},

		//Frontend requires the libs.js, js, ace, and gcss tasks, this is basically "grouping" tasks.
		slurp.Task{
			Name:  "frontend",
			Usage: "Run frontend tasks.",
			Deps:  []string{"assets", "index"},
			Action: func(c *slurp.C) error {
				return nil
			},
		},

		//The name says a lonet.
		slurp.Task{
			Name:  "watch",
			Usage: "Start watching gcss, ace, and javascript files and run crossponding tasks on change.",
			Deps:  []string{"frontend"},
			Action: func(c *slurp.C) error {

				g := watch.Watch(c, func(string) { b.Run(c, "index") }, "index.html")
				a := watch.Watch(c, func(string) { b.Run(c, "assets") }, "assets/*/*/*/*/*")

				<-c.Done()
				g.Close()
				a.Close()
				return nil
			},
		},

		//This will generate the resource file.
		slurp.Task{
			Name:  "embed",
			Usage: "compile public fileder into an http.FileSystem resource.",
			Action: func(c *slurp.C) error {
				return fs.Src(c,
					"public/*",
					"public/*/*",
				).Then(
					resources.Build(c, resources.Config{
						Pkg:     "main",
						Var:     "Public",
						Declare: false,
						Tag:     "embed",
					}),
					fs.Dest(c, "."),
				)
			},
		},

		//Start a livereload server and triggered everytime anything in public folder changes.
		slurp.Task{
			Name:  "livereload",
			Usage: "Start a tiny-lr server and monitor file changes in Pubic directory.",
			Action: func(c *slurp.C) error {

				l := watch.Watch(c, livereload.Start(c, config.Livereload, "public"),
					"public/*",
					"public/assets/*",
				)

				<-c.Done()
				l.Close()

				return nil
			},
		},

		// # Special tasks
		// when running this task with "slurp" it will run `go get`
		// for build dependenceis.
		slurp.Task{
			Name:  "init",
			Usage: "Assets.",
			Deps:  []string{"assets"},
			Action: func(c *slurp.C) error {
				//ideal for checking deps.
				return nil
			},
		},

		//When running slurp with no args, well, the "default" task is run.
		slurp.Task{
			Name:  "default",
			Usage: "Start livereload, watch, and gin tasks.",
			Deps:  []string{"livereload", "watch", "gin"},
			Action: func(c *slurp.C) error {
				//ideal for clean up.
				return nil
			},
		},
	)
}
