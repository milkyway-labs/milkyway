# Docs generator

This is a simple tool that can be used to generate the documentation that we can later on 
upload to GitBook.

This tool has been intended for reading the content of a markdown file and replace 
GitHub code reference sections with their content.  

For example if we want to include some protobuf code we can add a section like this:
````
```protobuf reference
https://github.com/cosmos/cosmos-sdk/blob/master/proto/cosmos/base/v1beta1/coin.proto#L12-L32
```
````
The tool will then get the content of `proto/cosmos/base/v1beta1/coin.proto` file from line 12 
to line 32 and write it to the resulting file.

## Setup

If this is your first time using this tool you can setup the python environment 
by running:

```bash
make install-requirements
```

## Generate documentation

The tools takes 2 arguments, the first is the path to the directory from which 
we generate the documentation and the second is the path to the directory where
the generated documentation will be saved.  

If you would like to generate the documentation for the modules inside the `x` directory
you can directly run the tool with the following command:

```bash
python docs-gen.py "./../../x" "./../../docs/generated/x
```

**Note**: This  tool will recursively explore the path that you provide and
stores the generated documentation in the path you provided preserving the original structure.  

### Generate modules documentation

To generate the documentation for our custom modules you can run:

```bash
make gen
```




