package ambari

import (
	"fmt"
	"log"
	"time"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/dghubble/sling"
)

type Privilege struct {
    PrivilegeInfo struct {
        PrivilegeId string `json:"privilege_id,omitempty"`,
        PermissionLabel string `json:"permission_label,omitempty"`,
        PermissionName string `json:"permission_name"`,
        PrincipalName string `json:"principal_name"`,
        PrincipalType string `json:"principal_type"`,
    } `json:"PrivilegeInfo"`,
}

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
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}


func resourceAmbariPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][ambari] Creating Privilege: %s", d.Privilege_id())
	clusterName := d.Get("cluster_name").(string)
	permissionName := d.Get("permission_name").(string)
	principalName := d.Get("principal_name").(string)
	principalType := d.Get("principal_type").(string)
	
	client := meta.(*Client).Client()
	path := fmt.Sprintf("/clusters/%s/privileges", clusterName)
	
	privilege := &Privilege {
        PrivilegeInfo.PermissionName: permissionName,
        PrivilegeInfo.PrincipalName: principalName,
        PrivilegeInfo.PrincipalType: principalType,
	}

    //Create the new privilege
	rep, err := client.Post(path).BodyJSON(privilege).Request()
	if err != nil {
		return err
	}
	
	// Get the privilege to get it's ID
	newPrivilege, err := getPrivilege(clusterName, permissionName, principalName, principalType, client)
	if err != nil {
	    return nil
	}
	if newPrivilege == nil {
	    return fmt.Errorf("Privilege ID not found")
	}
	
	d.SetId(newPrivilege.PrivilegeInfo.PrivilegeId)
	log.Printf("[INFO] Privilege ID: %s", d.Id()
	
	return resourceAmbariPrivilegeRead(d, meta)
}

func resourceAmbariPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Privilege: %s", d.Id())
	
	clusterName := d.Get("cluster_name").(string)
	privilegeId := d.Id()
	client := meta.(*Client).Client()
	path := fmt.Sprintf("/clusters/%s/privileges/%s", clusterName, privilegeId)
	
	privilege := new(Privilege)
    rep, err := client.GET(path).ReceiveSuccess(privilege)
    if err != nil {
        return err
    }
	if privilege == nil {
		return fmt.Errorf("Privilege not found with ID %s", privilegeId)
	}

	log.Printf("[INFO] Privilege Name: %s", privilege.PrivilegeInfo.PrincipalName)

	d.Set("permission_label", privilege.PrivilegeInfo.PermissionLabel)
	d.Set("permission_name", privilege.PrivilegeInfo.PermissionName)
	d.Set("principal_name", privilege.PrivilegeInfo.PrincipalName)
	d.Set("principal_type", privilege.PrivilegeInfo.PrincipalType)


	return nil
}

func resourceAmbariPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Privilege: %s", d.Id())
	clusterName := d.Get("cluster_name").(string)
	privilegeId := d.Id()
	permissionName := d.Get("permission_name").(string)
	principalName := d.Get("principal_name").(string)
	principalType := d.Get("principal_type").(string)
	client := meta.(*Client).Client()
	path := fmt.Sprintf("/clusters/%s/privileges/%s", clusterName, privilegeId)


	privilege := &Privilege {
        PrivilegeInfo.PermissionName: permissionName,
        PrivilegeInfo.PrincipalName: principalName,
        PrivilegeInfo.PrincipalType: principalType,
	}

    //Create the new privilege
	rep, err := client.Put(path).BodyJSON(privilege).Request()
	if err != nil {
		return err
	}

	return resourceAmbariPrivilegeRead(d, meta)
}

func resourceAmbariPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Privilege: %s", d.Id())
	clusterName := d.Get("cluster_name").(string)
	privilegeId := d.Id()
	client := meta.(*Client).Client()
	path := fmt.Sprintf("/clusters/%s/privileges/%s", clusterName, privilegeId)
	
	err := client.Delete(path).Request()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceAmbariPrivilegeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterName := d.Get("cluster_name").(string)
	privilegeId := d.Id()
	
	err := resourceAmbariPrivilegeRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}



// Permit to get privilege by search to obtain the ID
// Return the privilege if found
// Return nil if not found
func getPrivilege(string clusterName, string permissionName, string principalName, string principal_type, client *Sling ) (*Privilege, error) {
    
    if client == nil {
        fmt.Errorf("Client must be provided")
    }
    if clusterName == "" {
        fmt.Errorf("ClusterName must be provided")
    }
    
    path := fmt.Sprintf("/clusters/%s/privileges?PrivilegeInfo/permission_name=\"%s\"&PrivilegeInfo/principal_name=\"%s\"&PrivilegeInfo/principal_type=\"%s\"", clusterName, permissionName, principalName, principalType)
    privileges := new([]Privilege)
    rep, err := client.GET(path).ReceiveSuccess(privileges)
    if err != nil {
        return err
    }
    
    if len(privileges) == 1 {
        return &privileges[0]
    } else {
        return nil
    }
}

