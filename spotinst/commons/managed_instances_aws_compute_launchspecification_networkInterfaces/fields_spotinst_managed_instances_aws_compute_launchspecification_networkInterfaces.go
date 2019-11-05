package managed_instances_aws_compute_launchspecification_networkInterfaces

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
	"strconv"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[NetworkInterface] = commons.NewGenericField(
		commons.ManagedInstanceAwsNetworkInterfaces,
		NetworkInterface,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(DeviceIndex): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(AssociatePublicIpAddress): {
						Type:     schema.TypeBool,
						Optional: true,
					},

					string(AssociateIPV6Address): {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value []interface{} = nil
			if managedInstance.Compute != nil && managedInstance.Compute.LaunchSpecification != nil &&
				managedInstance.Compute.LaunchSpecification.NetworkInterfaces != nil {
				networkInterfaces := managedInstance.Compute.LaunchSpecification.NetworkInterfaces
				value = flattenAWSManagedInstanceNetworkInterfaces(networkInterfaces)
			}
			if value != nil {
				if err := resourceData.Set(string(NetworkInterface), value); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(NetworkInterface), err)
				}
			} else {
				if err := resourceData.Set(string(NetworkInterface), []*aws.NetworkInterface{}); err != nil {
					return fmt.Errorf(string(commons.FailureFieldReadPattern), string(NetworkInterface), err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(NetworkInterface)); ok {

				if interfaces, err := expandAWSManagedInstanceNetworkInterfaces(v); err != nil {
					return err
				} else {
					managedInstance.Compute.LaunchSpecification.SetNetworkInterfaces(interfaces)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value []*aws.NetworkInterface = nil
			if v, ok := resourceData.GetOk(string(NetworkInterface)); ok {
				if interfaces, err := expandAWSManagedInstanceNetworkInterfaces(v); err != nil {
					return err
				} else {
					value = interfaces
				}
			}
			managedInstance.Compute.LaunchSpecification.SetNetworkInterfaces(value)
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func flattenAWSManagedInstanceNetworkInterfaces(networkInterfaces []*aws.NetworkInterface) []interface{} {
	result := make([]interface{}, 0, len(networkInterfaces))
	for _, iface := range networkInterfaces {
		m := make(map[string]interface{})
		m[string(AssociatePublicIpAddress)] = spotinst.BoolValue(iface.AssociatePublicIPAddress)
		m[string(AssociateIPV6Address)] = spotinst.BoolValue(iface.AssociateIPV6Address)

		if iface.DeviceIndex != nil {
			m[string(DeviceIndex)] = strconv.Itoa(spotinst.IntValue(iface.DeviceIndex))
		}

		result = append(result, m)
	}
	return result
}

func expandAWSManagedInstanceNetworkInterfaces(data interface{}) ([]*aws.NetworkInterface, error) {
	list := data.(*schema.Set).List()
	interfaces := make([]*aws.NetworkInterface, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		networkInterface := &aws.NetworkInterface{}

		//todo sali - check this if
		if v, ok := m[string(NetworkInterfaceId)].(string); ok && v != "" {
			if v, ok := m[string(AssociatePublicIpAddress)].(bool); ok && v {
				return nil, errors.New("invalid Network interface: associate_public_ip_address must be undefined when using network_interface_id")
			}
			networkInterface.SetId(spotinst.String(v))
		} else {
			// AssociatePublicIp cannot be set at all when NetworkInterfaceId is specified
			if v, ok := m[string(AssociatePublicIpAddress)].(bool); ok {
				networkInterface.SetAssociatePublicIPAddress(spotinst.Bool(v))
			}
		}
		if v, ok := m[string(DeviceIndex)].(string); ok && v != "" {
			if intVal, err := strconv.Atoi(v); err != nil {
				return nil, err
			} else {
				networkInterface.SetDeviceIndex(spotinst.Int(intVal))
			}
		}

		if v, ok := m[string(AssociateIPV6Address)].(bool); ok {
			networkInterface.SetAssociateIPV6Address(spotinst.Bool(v))
		}

		interfaces = append(interfaces, networkInterface)
	}

	return interfaces, nil
}
