# Condense

CloudFormation template Preprocessor

 * Split templates into small "component" files, automatically included when referenced via `{"Ref": "ComponentName"}` or `{"Fn::GetAtt", ["ComponentName", "AttName"]}`
 * Automatic lookup of external stack Outputs via `{"Fn::GetAtt", ["StackName", "Outputs.OutputName"]}`
 * Locally-derefenced parameters via `-parameters` files, derefenced via `{"Ref": "ParameterKey"}` or `{"Fn::GetAtt", ["ParameterKey", "SubKey.SubSubKey"]}`
 * Add comments almost anywhere via JSON `"$comment"` keys


## Basic Example
```bash
condense --template=./stack.json --parameters=parameters.json
```
##### stack.json
```javascript
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
```javascript
{
  "VPCName": "aVPCName"
}
```
