package ambari

import (
	"fmt"
	restClient "github.com/disaster37/go-ambari-rest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func resourceAmbariPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAmbariPrivilegeCreate,
		Read:   resourceAmbariPrivilegeRead,
		Update: resourceAmbariPrivilegeUpdate,
		Delete: resourceAmbariPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAmbariPrivilegeImport,
		},

		Schema: map[string]*schema.Schema{
			"privilege_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_label": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CLUSTER.ADMINISTRATOR", "CLUSTER.OPERATOR", "SERVICE.ADMINISTRATOR", "SERVICE.OPERATOR", "CLUSTER.USER"}, true),
			},
			"principal_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GROUP", "USER"}, true),
			},
		},
	}
}

func resourceAmbariPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Info("Creating Privilege")
	clusterName := d.Get("cluster_name").(string)
	permissionName := d.Get("permission_name").(string)
	principalName := d.Get("principal_name").(string)
	principalType := d.Get("principal_type").(string)

	client := meta.(*Client).Client()

	privilege := &restClient.Privilege{
		PrivilegeInfo: &restClient.PrivilegeInfo{
			PermissionName: permissionName,
			PrincipalName:  principalName,
			PrincipalType:  principalType,
		},
	}
	log.Debugf("Privilege to create: %s", privilege.String())

	privilege, err := client.CreatePrivilege(clusterName, privilege)
	if err != nil {
		return err
	}
	log.Debugf("Privilege after to create: %s", privilege.String())

	d.SetId(strconv.FormatInt(privilege.PrivilegeInfo.PrivilegeId, 10))
	log.Infof("Privilege ID: %s", d.Id())

	return resourceAmbariPrivilegeRead(d, meta)
}

func resourceAmbariPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Refreshing Privilege: %s", d.Id())

	clusterName := d.Get("cluster_name").(string)
	privilegeId, _ := strconv.ParseInt(d.Id(), 10, 64)
	client := meta.(*Client).Client()
	log.Debugf("ClusterName: %s", clusterName)
	log.Debugf("PrivilegeId: %d", privilegeId)

	privilege, err := client.Privilege(clusterName, privilegeId)
	if err != nil {
		return err
	}

	if privilege == nil {
		return errors.New(fmt.Sprintf("Privilige with id %d not found", privilegeId))
	}

	log.Debug("Privilege after read: %s", privilege.String())
	log.Infof("Privilege ID : %s", privilege.PrivilegeInfo.PrincipalName)

	d.Set("privilege_id", d.Id())
	d.Set("permission_label", privilege.PrivilegeInfo.PermissionLabel)
	d.Set("permission_name", privilege.PrivilegeInfo.PermissionName)
	d.Set("principal_name", privilege.PrivilegeInfo.PrincipalName)
	d.Set("principal_type", privilege.PrivilegeInfo.PrincipalType)

	return nil
}

func resourceAmbariPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Updating Privilege: %s", d.Id())
	clusterName := d.Get("cluster_name").(string)
	privilegeId, _ := strconv.ParseInt(d.Id(), 10, 64)
	permissionName := d.Get("permission_name").(string)
	principalName := d.Get("principal_name").(string)
	principalType := d.Get("principal_type").(string)
	client := meta.(*Client).Client()

	// We do nothink if cluster name change
	privilege := &restClient.Privilege{
		PrivilegeInfo: &restClient.PrivilegeInfo{
			PermissionName: permissionName,
			PrincipalName:  principalName,
			PrincipalType:  principalType,
			PrivilegeId:    privilegeId,
		},
	}
	log.Debugf("Privilege to update: %s", privilege.String())

	privilege, err := client.UpdatePrivilege(clusterName, privilege)
	if err != nil {
		return err
	}
	log.Debugf("Privilege after update: %s", privilege.String())
	log.Infof("New Privilege ID: %d", privilege.PrivilegeInfo.PrivilegeId)

	d.SetId(strconv.FormatInt(privilege.PrivilegeInfo.PrivilegeId, 10))

	return resourceAmbariPrivilegeRead(d, meta)
}

func resourceAmbariPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Deleting Privilege: %s", d.Id())
	clusterName := d.Get("cluster_name").(string)
	privilegeId, _ := strconv.ParseInt(d.Id(), 10, 64)
	client := meta.(*Client).Client()

	log.Debugf("ClusterName: %s", clusterName)
	log.Debugf("PrivilegeId: %d", privilegeId)
	err := client.DeletePrivilege(clusterName, privilegeId)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceAmbariPrivilegeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	err := resourceAmbariPrivilegeRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
