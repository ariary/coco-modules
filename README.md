# 🌴 coco-modules

## Writing Module

*golang is not mandatory as long your respect IPC interaction and `coco` flow but is recommended*

Take inspiration of `./cmd/hello/hello.go`.

1. Write module

It has to implement:
* `Work(params, output)`:
    * `params`: parameters to run the module
    * `output`: used to send back data to `coco-agent`
    * It is the core of the module, put the payload there
* `GetName()`: return the name of the module
* `GetPrefix`: for output purpose

2. Build executable from module:
```golang
package main

func main() {
	module := YOURMODULESTRUCT{Name: "[MODULE_NAME]"}
    //connect to agent
	cc, err := module.ConnectToAgent("")
	if err != nil {
		fmt.Println("Failed to connect to agent:", err)
		os.Exit(92)
	}
	//wait for "connection" validation from agent
	module.WaitConnectionValidation(cc, module)
    //wait for instructions from agent
	module.WaitInstruction(cc, module)
}
```
Then compile it to fit the target requirmeents (OS, etc)

**Your module is ready**

## Available Modules:
* [x] hello
* [ ] wallpaper: change wallpaper
* [ ] desktop: write file on desktop
...
