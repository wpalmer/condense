# Condense

CloudFormation template Preprocessor

 * Allow external files to be included via `{"Fn::IncludeFile": "filename.json"}`
 * Automatic lookup of external stack Outputs via `{"Fn::GetAtt": ["StackName", "Outputs.OutputName"]}`
 * Locally-derefenced parameters via `-parameters` files, derefenced via `{"Ref": "ParameterKey"}` or `{"Fn::GetAtt": ["ParameterKey", "SubKey.SubSubKey"]}`
 * Add comments almost anywhere via JSON `"$comment"` keys
 * Use aliases to specify which reference to use, via `[...]`, eg: `{"Fn::GetAtt": ["[stacks.networks]", "Outputs.OutputName"]}`
 * Loops, via `{"Fn::For": ["$bindVar", ["data", "to", "iterate"], {"template": {"Ref": "$bindVar"}}]}`
 * more!

## Basic Example
```bash
condense --template=./stack.json --parameters=parameters.json
```
##### stack.json
```json
{
  "Resources": {
    "VPC": {
      "$comment": "The primary VPC",
      "Type": "AWS::EC2::VPC",
      "Properties": {
        "CidrBlock": "10.0.0.0/16",
        "Tags": [{"Key": "Name", "Value": {"Ref": "VPCName"}}]
      }
    }
  }
}
```
##### parameters.json
```json
{
  "VPCName": "aVPCName"
}
```

## Rules

The template preprocessor visits each node in the template, passing each
visited node through the set of attached "Rules". Most rules are
analogous to "functions", returning a result based on the contents of
the node. Most rules also "fall-through" when they don't match, assuming
that they will be processed (or detected as errors) by CloudFormation.

### ExcludeComments

Removes object entries with the key "$comment", or entire objects
containing *only* the key "$comment". eg:
```json
{
    "a": "one",
    "$comment": "this will be removed",
    "b": [
      "b.a",
      {"$comment": "this whole object will be removed"},
      "b.b"
    ]
}
```

### FnAdd

Adds the floating-point values within the supplied array. Mostly
intended for tweaking array indeces or other derived values, eg:
```json
{"Fn::Add": [1, 3]}
```
Outputs:
```json
4
```

### FnIf, FnEquals, FnAnd, FnOr, FnNot

Analogous to CloudFormation's `Fn::If`, `Fn::Equals`, `Fn::And`,
`Fn::Or`, and `Fn::Not`. Allows these rules to be processed early, to
reduce final template size.

### FnConcat

Combine arrays into a single array, eg:
```json
{"Fn::Concat": [[1,2], [3,4]]}
```
Outputs:
```json
[1,2,3,4]
```

### FnFindFile

Search for a particular file among several candidate directories.
Intended to allow Included files to be optionally overridden for
testing. eg:
```json
{"Fn::FindFile": ["local", "default"], "included.json"}
```

### FnFor

Iterate over a list of values, applying each value to the specified
template, eg:
```json
{"Fn::For": [
  ["$i", "$value"],
  ["a", "b", "c"],
  {
    "index": {"Ref": "$i"},
    "value": {"Ref": "$value"}
  }
]}
```
Outputs:
```json
[
  {"index": 0, "value": "a"},
  {"index": 1, "value": "b"},
  {"index": 2, "value": "c"}
]
```

### FnFromEntries

Convert a list of `{"key": ..., "value": ...}` pairs to a single object.
Combined with Fn::For, new objects can be built via iteration. eg:
```json
{"Fn::FromEntries": [
  {"key": "a", "value": "one"},
  {"key": "b", "value": "two"}
]}
```
Outputs:
```json
{"a": "one", "b": "two"}
```

### FnGetAtt, Ref

Analogous to the CloudFormation `Fn::GetAtt` and `Ref` functions, with
the added abilities of being able to reference locally-provided
parameters, external CloudFormation stacks, and bound variables from
functions such as `Fn::For` and `Fn::With`

### FnIncludeFile

Reference an external file, adding it (as its JSON interpretation) to the
emplate. This function will panic if the file is not found, or is not valid
JSON.

### FnJoin

Analogous to the CloudFormation `Fn::Join` method, but allowing for
early processing in order to reduce final template size, eg:
```json
{"Fn::Join": [",", ["a", "b", "c"]]}
```
Outputs:
```json
"a,b,c"
```

### FnKeys

Return the keys of an object, as an array. The order of returned keys is
not stable. eg:
```json
{"Fn::Keys": {"a": "one", "b": "two"}}
```
Outputs:
```json
["b", "a"]
```

### FnMerge

Perform a simple merge of the upper-most keys of an object, eg:
```json
{"Fn::Merge": [{"a": "one"}, {"b": "two"}]}
```
Outputs:
```json
{"a": "one", "b": "two"}
```

### FnMergeDeep

Perform a deep merge of an object, to a specified depth, eg:
```json
{"Fn::MergeDeep": [1, [
  {
    "a": "one.a",
    "b": "one.b",
    "deep": {
      "deep.a": "one.deep.a",
      "deep.b": "one.deep.b",
      "deep.deep": {
        "deep.deep.a": "one.deep.deep.a",
        "deep.deep.b": "one.deep.deep.b"
      }
    }
  },
  {
    "b": "two.b",
    "c": "two.c",
    "deep": {
      "deep.b": "two.deep.b",
      "deep.c": "two.deep.c",
      "deep.deep": {
        "deep.deep.b": "one.deep.deep.b",
        "deep.deep.c": "one.deep.deep.c"
      }
    }
  },
  {
    "d": "three.d",
    "deep": {
      "deep.d": "three.deep.d",
      "deep.deep": {
        "deep.deep.d": "three.deep.deep.d"
      }
    }
  }
]]}
```
Outputs:
```json
{
  "a": "one.a",
  "b": "two.b",
  "c": "two.c",
  "d": "three.d",
  "deep": {
    "deep.a": "one.deep.a",
    "deep.b": "two.deep.b",
    "deep.c": "two.deep.c",
    "deep.d": "three.deep.d",
    "deep.deep": {
      "deep.deep.d": "three.deep.deep.d"
    }
  }
}
```

### FnSplit

The inverse of `Fn::Join`, converts a string to an array, eg:
```json
{"Fn::Split": [",", "a,b,c"]}
```
Outputs:
```json
["a", "b", "c"]
```

### FnToEntries

The inverse of `Fn::FromEntries`. As with `Fn::Keys`, the order of the
resulting array is not stable. eg:
```json
{"Fn::ToEntries": {"a": "one", "b": "two"}}
```
Outputs:
```json
[
  {"key": "b", "value": "two"},
  {"key": "a", "value": "one"}
]
```

### FnUnique

Removes duplicate values from an array, eg:
```json
{"Fn::Unique": ["a", "b", "c", "a", "c"]}
```
Outputs:
```json
["a", "b", "c"]
```

### FnWith

Passes bound values into a specified template. Usually used when the
template is an Included file, eg:
```json
{"Fn::With": [{"$a": "foo"}, {"Ref": "$a"}]}
```
Outputs:
```json
"foo"
```

### ReduceConditions

Locates nodes within the `Conditions` section of the template, and
ensures they are converted to boolean equivalents, rather than true
booleans, which CloudFormation can't process. eg:
```json
{"Conditions": {"aCondition": false}}
```
Becomes:
```json
{"Conditions": {"aCondition": {"Fn::Equals": ["1", "0"]}}}
```
