# MacOS File Magic Number Detection

Observes the directory chosen for files that their advertised file types do not match their true file type based on the magic number. This then sends a notification to the user alerting them of the file.




# Usage

```
go run src/main.go

```

```
-filepath string
        Please enter a target directory (default "./")
  -logger string
        Please specify the log file location (default "./")
```


# Dependancies

github.com/gen2brain/beeep
github.com/go-toast/toast
github.com/godbus/dbus/v5
github.com/nu7hatch/gouuid
github.com/tadvi/systray

golang.org/x/exp
golang.org/x/sys