[DEMO LINK](https://therecipe.github.io/widgets_playground)

This is a showcase example for the "js" and "wasm" targets and also the new JavaScript API of [therecipe/qt](https://github.com/therecipe/qt)

General installation instructions for `therecipe/qt` can be found here: https://github.com/therecipe/qt/wiki/Installation

If you are already familiar with qtdeploy and the docker deployments then just pull the "js" or "wasm" image and deploy with `qtdeploy -docker build js` or "wasm" as usual.

---

If you are new and want to build this yourself, then just take the following steps:

-	Install Go: https://golang.org/dl/

-	Install Git: https://git-scm.com/downloads

-	Install tooling: `go install -v github.com/therecipe/qt/cmd/...`

-	Pull the repo: `go get -v -d github.com/therecipe/widgets_playground`

-	Install Docker: https://store.docker.com/search?offering=community&type=edition

	-	On Windows: [share](https://docs.docker.com/docker-for-windows/#shared-drives) the drive containing your **GOPATH** with docker
	-	On Linux: if necessary run docker as [root](https://docs.docker.com/install/linux/linux-postinstall/#manage-docker-as-a-non-root-user)
	-	On macOS: [share](https://docs.docker.com/docker-for-mac/#file-sharing) your **GOPATH** with docker if it isn't located in some subfolder below `/Users/`, `/Volumes/`, `/private/` or `/tmp/`

-	Pull the "js" or "wasm" image:

	```
	docker pull therecipe/qt:js
	```

	or

	```
	docker pull therecipe/qt:wasm
	```

-	Run the deployment: (replace "js" with "wasm" for an full WebAssembly build)

	```
	cd $(go env GOPATH)/src/github.com/therecipe/widgets_playground
	$(go env GOPATH)/bin/qtdeploy -docker build js
	```

-	You should find the deployed application in the `deploy/js` or `deploy/wasm` subdir

-	Open `deploy/{js|wasm}/index.html` with your browser

---

You can most of the time use `qtdeploy -docker -fast build js` after you made some minor changes, heavier changes like the introduction of new Qt classes or functions, or changes inside the JavaScript files will force you to compile without the "-fast" flag once again.

It's planned to remove these limitations in the future.
