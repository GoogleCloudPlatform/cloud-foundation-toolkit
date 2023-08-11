# Module Swapper

Module Swapper is a utility used for swapping TF registry references with local modules. It will ignore registry references to all other modules except for the one in current directory.

```
Usage of module-swapper:
  -examples-path string
        Path to examples that should be swapped. Defaults to cwd/examples (default "examples")
  -registry-prefix string
        Module registry prefix (default "terraform-google-modules")
  -registry-suffix string
        Module registry suffix (default "google")
  -restore
        Restores disabled modules
  -submods-path string
        Path to a submodules if any that maybe referenced. Defaults to working dir/modules (default "modules")
  -workdir string
        Absolute path to root module where examples should be swapped. Defaults to working directory
```