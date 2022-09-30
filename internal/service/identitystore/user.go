package identitystore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/document"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceUserCreate,
		ReadWithoutTimeout:   resourceUserRead,
		UpdateWithoutTimeout: resourceUserUpdate,
		DeleteWithoutTimeout: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"formatted": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"locality": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"postal_code": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"region": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"street_address": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
					},
				},
			},
			"display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"emails": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
					},
				},
			},
			"external_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"locale": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"name": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"formatted": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"given_name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"honorific_prefix": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"honorific_suffix": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"middle_name": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
					},
				},
			},
			"nickname": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"phone_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
						"value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
						},
					},
				},
			},
			"preferred_language": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"profile_url": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"timezone": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"title": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 128)),
			},
			"user_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 1024)),
			},
		},
	}
}

const (
	ResNameUser = "User"
)

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).IdentityStoreConn

	in := &identitystore.CreateUserInput{
		DisplayName:     aws.String(d.Get("display_name").(string)),
		IdentityStoreId: aws.String(d.Get("identity_store_id").(string)),
		UserName:        aws.String(d.Get("user_name").(string)),
	}

	if v, ok := d.GetOk("addresses"); ok && len(v.([]interface{})) > 0 {
		in.Addresses = expandAddresses(v.([]interface{}))
	}

	if v, ok := d.GetOk("emails"); ok && len(v.([]interface{})) > 0 {
		in.Emails = expandEmails(v.([]interface{}))
	}

	if v, ok := d.GetOk("locale"); ok {
		in.Locale = aws.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		in.Name = expandName(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("phone_numbers"); ok && len(v.([]interface{})) > 0 {
		in.PhoneNumbers = expandPhoneNumbers(v.([]interface{}))
	}

	if v, ok := d.GetOk("nickname"); ok {
		in.NickName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_language"); ok {
		in.PreferredLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("profile_url"); ok {
		in.ProfileUrl = aws.String(v.(string))
	}

	if v, ok := d.GetOk("timezone"); ok {
		in.Timezone = aws.String(v.(string))
	}

	if v, ok := d.GetOk("title"); ok {
		in.Title = aws.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		in.UserType = aws.String(v.(string))
	}

	out, err := conn.CreateUser(ctx, in)
	if err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionCreating, ResNameUser, d.Get("identity_store_id").(string), err)
	}

	if out == nil || out.UserId == nil {
		return create.DiagError(names.IdentityStore, create.ErrActionCreating, ResNameUser, d.Get("identity_store_id").(string), errors.New("empty output"))
	}

	d.SetId(fmt.Sprintf("%s/%s", aws.ToString(out.IdentityStoreId), aws.ToString(out.UserId)))

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).IdentityStoreConn

	identityStoreId, userId, err := resourceUserParseID(d.Id())

	if err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionReading, ResNameUser, d.Id(), err)
	}

	out, err := findUserByID(ctx, conn, identityStoreId, userId)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IdentityStore User (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionReading, ResNameUser, d.Id(), err)
	}

	d.Set("display_name", out.DisplayName)
	d.Set("identity_store_id", out.IdentityStoreId)
	d.Set("locale", out.Locale)
	d.Set("nickname", out.NickName)
	d.Set("preferred_language", out.PreferredLanguage)
	d.Set("profile_url", out.ProfileUrl)
	d.Set("timezone", out.Timezone)
	d.Set("title", out.Title)
	d.Set("user_id", out.UserId)
	d.Set("user_name", out.UserName)
	d.Set("user_type", out.UserType)

	if err := d.Set("addresses", flattenAddresses(out.Addresses)); err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionSetting, ResNameUser, d.Id(), err)
	}

	if err := d.Set("emails", flattenEmails(out.Emails)); err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionSetting, ResNameUser, d.Id(), err)
	}

	if err := d.Set("external_ids", flattenExternalIds(out.ExternalIds)); err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionSetting, ResNameUser, d.Id(), err)
	}

	if err := d.Set("name", []interface{}{flattenName(out.Name)}); err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionSetting, ResNameUser, d.Id(), err)
	}

	if err := d.Set("phone_numbers", flattenPhoneNumbers(out.PhoneNumbers)); err != nil {
		return create.DiagError(names.IdentityStore, create.ErrActionSetting, ResNameUser, d.Id(), err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).IdentityStoreConn

	in := &identitystore.UpdateUserInput{
		IdentityStoreId: aws.String(d.Get("identity_store_id").(string)),
		UserId:          aws.String(d.Get("user_id").(string)),
		Operations:      nil,
	}

	// IMPLEMENTATION NOTE.
	//
	// Complex types, such as the `emails` field, don't allow field by field
	// updates, and require that the entire sub-object is modified.
	//
	// In those sub-objects, to remove a field, it must not be present at all
	// in the updated attribute value.
	//
	// However, structs such as types.Email don't specify omitempty in their
	// struct tags, so the document.NewLazyDocument marshaller will write out
	// nulls.
	//
	// This is why, for those complex fields, a custom Expand function is
	// provided that converts the Go SDK type (e.g. types.Email) into a field
	// by field representation of what the API would expect.

	fieldsToUpdate := []struct {
		// Attribute corresponds to the provider schema.
		Attribute string

		// Field corresponds to the AWS API schema.
		Field string

		// Expand, when not nil, is used to transform the value of the field
		// given in Attribute before it's passed to the UpdateOperation.
		Expand func(interface{}) interface{}
	}{
		{
			Attribute: "display_name",
			Field:     "displayName",
		},
		{
			Attribute: "locale",
			Field:     "locale",
		},
		{
			Attribute: "name.0.family_name",
			Field:     "name.familyName",
		},
		{
			Attribute: "name.0.formatted",
			Field:     "name.formatted",
		},
		{
			Attribute: "name.0.given_name",
			Field:     "name.givenName",
		},
		{
			Attribute: "name.0.honorific_prefix",
			Field:     "name.honorificPrefix",
		},
		{
			Attribute: "name.0.honorific_suffix",
			Field:     "name.honorificSuffix",
		},
		{
			Attribute: "name.0.middle_name",
			Field:     "name.middleName",
		},
		{
			Attribute: "nickname",
			Field:     "nickName",
		},
		{
			Attribute: "preferred_language",
			Field:     "preferredLanguage",
		},
		{
			Attribute: "profile_url",
			Field:     "profileUrl",
		},
		{
			Attribute: "timezone",
			Field:     "timezone",
		},
		{
			Attribute: "title",
			Field:     "title",
		},
		{
			Attribute: "user_type",
			Field:     "userType",
		},
		{
			Attribute: "addresses",
			Field:     "addresses",
			Expand: func(value interface{}) interface{} {
				addresses := expandAddresses(value.([]interface{}))

				var result []interface{}

				// The API requires a null to unset the list, so in the case
				// of no addresses, a nil result is preferable.
				for _, address := range addresses {
					m := map[string]interface{}{}

					if v := address.Country; v != nil {
						m["country"] = v
					}

					if v := address.Formatted; v != nil {
						m["formatted"] = v
					}

					if v := address.Locality; v != nil {
						m["locality"] = v
					}

					if v := address.PostalCode; v != nil {
						m["postalCode"] = v
					}

					m["primary"] = address.Primary

					if v := address.Region; v != nil {
						m["region"] = v
					}

					if v := address.StreetAddress; v != nil {
						m["streetAddress"] = v
					}

					if v := address.Type; v != nil {
						m["type"] = v
					}

					result = append(result, m)
				}

				return result
			},
		},
		{
			Attribute: "emails",
			Field:     "emails",
			Expand: func(value interface{}) interface{} {
				emails := expandEmails(value.([]interface{}))

				var result []interface{}

				// The API requires a null to unset the list, so in the case
				// of no emails, a nil result is preferable.
				for _, email := range emails {
					m := map[string]interface{}{}

					m["primary"] = email.Primary

					if v := email.Type; v != nil {
						m["type"] = v
					}

					if v := email.Value; v != nil {
						m["value"] = v
					}

					result = append(result, m)
				}

				return result
			},
		},
		{
			Attribute: "phone_numbers",
			Field:     "phoneNumbers",
			Expand: func(value interface{}) interface{} {
				emails := expandPhoneNumbers(value.([]interface{}))

				var result []interface{}

				// The API requires a null to unset the list, so in the case
				// of no emails, a nil result is preferable.
				for _, email := range emails {
					m := map[string]interface{}{}

					m["primary"] = email.Primary

					if v := email.Type; v != nil {
						m["type"] = v
					}

					if v := email.Value; v != nil {
						m["value"] = v
					}

					result = append(result, m)
				}

				return result
			},
		},
	}

	for _, fieldToUpdate := range fieldsToUpdate {
		if d.HasChange(fieldToUpdate.Attribute) {
			value := d.Get(fieldToUpdate.Attribute)

			if expand := fieldToUpdate.Expand; expand != nil {
				value = expand(value)
			}

			// The API doesn't allow empty attribute values. To unset an
			// attribute, set it to null.
			if value == "" {
				value = nil
			}

			in.Operations = append(in.Operations, types.AttributeOperation{
				AttributePath:  aws.String(fieldToUpdate.Field),
				AttributeValue: document.NewLazyDocument(value),
			})
		}
	}

	if len(in.Operations) > 0 {
		log.Printf("[DEBUG] Updating IdentityStore User (%s): %#v", d.Id(), in)
		_, err := conn.UpdateUser(ctx, in)
		if err != nil {
			return create.DiagError(names.IdentityStore, create.ErrActionUpdating, ResNameUser, d.Id(), err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).IdentityStoreConn

	log.Printf("[INFO] Deleting IdentityStore User %s", d.Id())

	_, err := conn.DeleteUser(ctx, &identitystore.DeleteUserInput{
		IdentityStoreId: aws.String(d.Get("identity_store_id").(string)),
		UserId:          aws.String(d.Get("user_id").(string)),
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.IdentityStore, create.ErrActionDeleting, ResNameUser, d.Id(), err)
	}

	return nil
}

func findUserByID(ctx context.Context, conn *identitystore.Client, identityStoreId, userId string) (*identitystore.DescribeUserOutput, error) {
	in := &identitystore.DescribeUserInput{
		IdentityStoreId: aws.String(identityStoreId),
		UserId:          aws.String(userId),
	}

	out, err := conn.DescribeUser(ctx, in)

	if err != nil {
		var e *types.ResourceNotFoundException
		if errors.As(err, &e) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		} else {
			return nil, err
		}
	}

	if out == nil || out.UserId == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func flattenAddress(apiObject *types.Address) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.Country; v != nil {
		m["country"] = aws.ToString(v)
	}

	if v := apiObject.Formatted; v != nil {
		m["formatted"] = aws.ToString(v)
	}

	if v := apiObject.Locality; v != nil {
		m["locality"] = aws.ToString(v)
	}

	if v := apiObject.PostalCode; v != nil {
		m["postal_code"] = aws.ToString(v)
	}

	m["primary"] = apiObject.Primary

	if v := apiObject.Region; v != nil {
		m["region"] = aws.ToString(v)
	}

	if v := apiObject.StreetAddress; v != nil {
		m["street_address"] = aws.ToString(v)
	}

	if v := apiObject.Type; v != nil {
		m["type"] = aws.ToString(v)
	}

	return m
}

func expandAddress(tfMap map[string]interface{}) *types.Address {
	if tfMap == nil {
		return nil
	}

	a := &types.Address{}

	if v, ok := tfMap["country"].(string); ok && v != "" {
		a.Country = aws.String(v)
	}

	if v, ok := tfMap["formatted"].(string); ok && v != "" {
		a.Formatted = aws.String(v)
	}

	if v, ok := tfMap["locality"].(string); ok && v != "" {
		a.Locality = aws.String(v)
	}

	if v, ok := tfMap["postal_code"].(string); ok && v != "" {
		a.PostalCode = aws.String(v)
	}

	a.Primary = tfMap["primary"].(bool)

	if v, ok := tfMap["region"].(string); ok && v != "" {
		a.Region = aws.String(v)
	}

	if v, ok := tfMap["street_address"].(string); ok && v != "" {
		a.StreetAddress = aws.String(v)
	}

	if v, ok := tfMap["type"].(string); ok && v != "" {
		a.Type = aws.String(v)
	}

	return a
}

func flattenAddresses(apiObjects []types.Address) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		apiObject := apiObject
		l = append(l, flattenAddress(&apiObject))
	}

	return l
}

func expandAddresses(tfList []interface{}) []types.Address {
	s := make([]types.Address, 0, len(tfList))

	for _, r := range tfList {
		m, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		a := expandAddress(m)

		if a == nil {
			continue
		}

		s = append(s, *a)
	}

	return s
}

func flattenName(apiObject *types.Name) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.FamilyName; v != nil {
		m["family_name"] = aws.ToString(v)
	}

	if v := apiObject.Formatted; v != nil {
		m["formatted"] = aws.ToString(v)
	}

	if v := apiObject.GivenName; v != nil {
		m["given_name"] = aws.ToString(v)
	}

	if v := apiObject.HonorificPrefix; v != nil {
		m["honorific_prefix"] = aws.ToString(v)
	}

	if v := apiObject.HonorificSuffix; v != nil {
		m["honorific_suffix"] = aws.ToString(v)
	}

	if v := apiObject.MiddleName; v != nil {
		m["middle_name"] = aws.ToString(v)
	}

	return m
}

func expandName(tfMap map[string]interface{}) *types.Name {
	if tfMap == nil {
		return nil
	}

	a := &types.Name{}

	if v, ok := tfMap["family_name"].(string); ok && v != "" {
		a.FamilyName = aws.String(v)
	}

	if v, ok := tfMap["formatted"].(string); ok && v != "" {
		a.Formatted = aws.String(v)
	}

	if v, ok := tfMap["given_name"].(string); ok && v != "" {
		a.GivenName = aws.String(v)
	}

	if v, ok := tfMap["honorific_prefix"].(string); ok && v != "" {
		a.HonorificPrefix = aws.String(v)
	}

	if v, ok := tfMap["honorific_suffix"].(string); ok && v != "" {
		a.HonorificSuffix = aws.String(v)
	}

	if v, ok := tfMap["middle_name"].(string); ok && v != "" {
		a.MiddleName = aws.String(v)
	}

	return a
}

func flattenEmail(apiObject *types.Email) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	m["primary"] = apiObject.Primary

	if v := apiObject.Type; v != nil {
		m["type"] = aws.ToString(v)
	}

	if v := apiObject.Value; v != nil {
		m["value"] = aws.ToString(v)
	}

	return m
}

func expandEmail(tfMap map[string]interface{}) *types.Email {
	if tfMap == nil {
		return nil
	}

	a := &types.Email{}

	a.Primary = tfMap["primary"].(bool)

	if v, ok := tfMap["type"].(string); ok && v != "" {
		a.Type = aws.String(v)
	}

	if v, ok := tfMap["value"].(string); ok && v != "" {
		a.Value = aws.String(v)
	}

	return a
}

func flattenEmails(apiObjects []types.Email) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		apiObject := apiObject
		l = append(l, flattenEmail(&apiObject))
	}

	return l
}

func expandEmails(tfList []interface{}) []types.Email {
	s := make([]types.Email, 0, len(tfList))

	for _, r := range tfList {
		m, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		a := expandEmail(m)

		if a == nil {
			continue
		}

		s = append(s, *a)
	}

	return s
}

func flattenExternalId(apiObject *types.ExternalId) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	if v := apiObject.Id; v != nil {
		m["id"] = aws.ToString(v)
	}

	if v := apiObject.Issuer; v != nil {
		m["issuer"] = aws.ToString(v)
	}

	return m
}

func flattenExternalIds(apiObjects []types.ExternalId) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		apiObject := apiObject
		l = append(l, flattenExternalId(&apiObject))
	}

	return l
}

func flattenPhoneNumber(apiObject *types.PhoneNumber) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{}

	m["primary"] = apiObject.Primary

	if v := apiObject.Type; v != nil {
		m["type"] = aws.ToString(v)
	}

	if v := apiObject.Value; v != nil {
		m["value"] = aws.ToString(v)
	}

	return m
}

func expandPhoneNumber(tfMap map[string]interface{}) *types.PhoneNumber {
	if tfMap == nil {
		return nil
	}

	a := &types.PhoneNumber{}

	a.Primary = tfMap["primary"].(bool)

	if v, ok := tfMap["type"].(string); ok && v != "" {
		a.Type = aws.String(v)
	}

	if v, ok := tfMap["value"].(string); ok && v != "" {
		a.Value = aws.String(v)
	}

	return a
}

func flattenPhoneNumbers(apiObjects []types.PhoneNumber) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var l []interface{}

	for _, apiObject := range apiObjects {
		apiObject := apiObject
		l = append(l, flattenPhoneNumber(&apiObject))
	}

	return l
}

func expandPhoneNumbers(tfList []interface{}) []types.PhoneNumber {
	s := make([]types.PhoneNumber, 0, len(tfList))

	for _, r := range tfList {
		m, ok := r.(map[string]interface{})

		if !ok {
			continue
		}

		a := expandPhoneNumber(m)

		if a == nil {
			continue
		}

		s = append(s, *a)
	}

	return s
}

func resourceUserParseID(id string) (identityStoreId, userId string, err error) {
	parts := strings.Split(id, "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		err = errors.New("expected a resource id in the form: identity-store-id/user-id")
		return
	}

	return parts[0], parts[1], nil
}