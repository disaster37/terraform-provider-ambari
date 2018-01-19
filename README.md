# terraform-provider-ambari
Ambari provider for Terraform


## Resources

### Privilege

Sample:

```
resource ambari_privilege "admins" {
  cluster_name = "hdp-test"
  permission_name = "CLUSTER.ADMINISTRATOR"
  principal_name = "admins"
  principal_type = "GROUP"
}
```

Parameters:
- **permission_name**: Can be CLUSTER.ADMINISTRATOR, CLUSTER.OPERATOR, SERVICE.ADMINISTRATOR, SERVICE.OPERATOR, CLUSTER.USER
- **principal_type**: Can be USER or GROUP
- **principal_name**: The user or group you should to affect the permission. The user/group must be already exist in Ambari.

