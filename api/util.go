package api

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// Resource types
	RESOURCE_GROUP  = "group"
	RESOURCE_USER   = "user"
	RESOURCE_POLICY = "policy"

	// Constraints
	MAX_EXTERNAL_ID_LENGTH = 128
	MAX_NAME_LENGTH        = 128
	MAX_ACTION_LENGTH      = 128
	MAX_PATH_LENGTH        = 512

	// Actions

	// User actions
	USER_ACTION_CREATE_USER          = "iam:CreateUser"
	USER_ACTION_DELETE_USER          = "iam:DeleteUser"
	USER_ACTION_GET_USER             = "iam:GetUser"
	USER_ACTION_LIST_USERS           = "iam:ListUsers"
	USER_ACTION_UPDATE_USER          = "iam:UpdateUser"
	USER_ACTION_LIST_GROUPS_FOR_USER = "iam:ListGroupsForUser"
	USER_ACTION_LIST_ORG_USERS       = "iam:ListOrgUsers"

	// Group actions
	GROUP_ACTION_CREATE_GROUP                 = "iam:CreateGroup"
	GROUP_ACTION_DELETE_GROUP                 = "iam:DeleteGroup"
	GROUP_ACTION_GET_GROUP                    = "iam:GetGroup"
	GROUP_ACTION_LIST_GROUPS                  = "iam:ListGroups"
	GROUP_ACTION_UPDATE_GROUP                 = "iam:UpdateGroup"
	GROUP_ACTION_LIST_MEMBERS                 = "iam:ListMembers"
	GROUP_ACTION_ADD_MEMBER                   = "iam:AddMember"
	GROUP_ACTION_REMOVE_MEMBER                = "iam:RemoveMember"
	GROUP_ACTION_ATTACH_GROUP_POLICY          = "iam:AttachGroupPolicy"
	GROUP_ACTION_DETACH_GROUP_POLICY          = "iam:DetachGroupPolicy"
	GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES = "iam:ListAttachedGroupPolicies"
	GROUP_ACTION_LIST_ALL_GROUPS              = "iam:ListAllGroups"

	// Policy actions
	POLICY_ACTION_CREATE_POLICY        = "iam:CreatePolicy"
	POLICY_ACTION_DELETE_POLICY        = "iam:DeletePolicy"
	POLICY_ACTION_UPDATE_POLICY        = "iam:UpdatePolicy"
	POLICY_ACTION_GET_POLICY           = "iam:GetPolicy"
	POLICY_ACTION_LIST_ATTACHED_GROUPS = "iam:ListAttachedGroups"
	POLICY_ACTION_LIST_POLICIES        = "iam:ListPolicies"
	POLICY_ACTION_LIST_ALL_POLICIES    = "iam:ListAllPolicies"
)

func CreateUrn(org string, resource string, path string, name string) string {
	switch resource {
	case RESOURCE_USER:
		return fmt.Sprintf("urn:iws:iam::user%v%v", path, name)
	default:
		return fmt.Sprintf("urn:iws:iam:%v:%v%v%v", org, resource, path, name)
	}
}

func GetUrnPrefix(org string, resource string, path string) string {
	switch resource {
	case RESOURCE_USER:
		return fmt.Sprintf("urn:iws:iam::user%v*", path)
	default:
		return fmt.Sprintf("urn:iws:iam:%v:%v%v*", org, resource, path)
	}
}

func IsValidUserExternalID(externalID string) bool {
	r, _ := regexp.Compile(`^[\w+.@\-]+$`)
	return r.MatchString(externalID) && len(externalID) < MAX_EXTERNAL_ID_LENGTH
}

// this func validates group and policy names
func IsValidName(name string) bool {
	r, _ := regexp.Compile(`^[\w\-]+$`)
	return r.MatchString(name) && len(name) < MAX_NAME_LENGTH
}

func IsValidPath(path string) bool {
	r, _ := regexp.Compile(`^/$|^/[\w+/\-]+\w+/$`)
	r2, _ := regexp.Compile(`[/]{2,}`)
	return r.MatchString(path) && !r2.MatchString(path) && len(path) < MAX_PATH_LENGTH
}

func IsValidEffect(effect string) error {
	if effect != "allow" && effect != "deny" {
		return &Error{
			Code:    REGEX_NO_MATCH,
			Message: fmt.Sprintf("No regex match in effect: %v", effect),
		}
	}
	return nil
}

func IsValidAction(actions []string) error {
	r, _ := regexp.Compile(`^[\w\-:]+[\w-*]+$`)
	r2, _ := regexp.Compile(`[*]{2,}|[:]{2,}`)
	for _, action := range actions {
		if !r.MatchString(action) || r2.MatchString(action) || len(action) > MAX_ACTION_LENGTH {
			return &Error{
				Code:    REGEX_NO_MATCH,
				Message: fmt.Sprintf("No regex match in action: %v", action),
			}
		}
	}
	return nil
}

func IsValidResources(resources []string) error {
	//err generator helper
	errFunc := func(resource string) error {
		return &Error{
			Code:    REGEX_NO_MATCH,
			Message: fmt.Sprintf("No regex match in resource: %v", resource),
		}
	}

	wordRegex, _ := regexp.Compile(`^[\w+\-.@]+$`)
	wordPrefixRegex, _ := regexp.Compile(`^[\w+\-.@]+\*$`)

	r, _ := regexp.Compile(`^\*$|^[\w+\-@.]+\*?$|^[\w+\-@.]+\*?$|^[\w+\-@.]+(/?([\w+\-@.]+/)*([\w+\-@.]|[*])+)?$`)
	r2, _ := regexp.Compile(`[/]{2,}|[:]{2,}|[*]{2,}`)

	for _, resource := range resources {
		blocks := strings.Split(resource, ":")
		for n, block := range blocks {
			switch n {
			case 0:
				if len(blocks) < 2 { // This is the last block
					if block != "*" {
						return errFunc(resource)
					}
				} else {
					if block != "urn" {
						return errFunc(resource)
					}
				}
			case 1:
				if len(blocks) < 3 { // This is the last block
					if block != "*" && !wordPrefixRegex.MatchString(block) {
						return errFunc(resource)
					}
				} else {
					if !wordRegex.MatchString(block) {
						return errFunc(resource)
					}
				}
			case 2:
				if len(blocks) < 4 { // This is the last block
					if block != "*" && !wordPrefixRegex.MatchString(block) {
						return errFunc(resource)
					}
				} else {
					if !wordRegex.MatchString(block) {
						return errFunc(resource)
					}
				}
			case 3:
				if len(blocks) < 5 { // This is the last block
					if block != "*" && !wordPrefixRegex.MatchString(block) {
						return errFunc(resource)
					}
				} else {
					if block != "" && !wordRegex.MatchString(block) {
						return errFunc(resource)
					}
				}
			case 4:
				if !r.MatchString(block) || r2.MatchString(block) {
					return errFunc(resource)
				}
			default:
				return &Error{
					Code:    INVALID_PARAMETER_ERROR,
					Message: fmt.Sprintf("Invalid resource definition: %v", resource),
				}

			}
		}
	}
	return nil
}

func IsValidStatement(statements *[]Statement) error {
	for _, statement := range *statements {
		err := IsValidEffect(statement.Effect)
		if err != nil {
			return err
		}
		err = IsValidAction(statement.Action)
		if err != nil {
			return err
		}
		err = IsValidResources(statement.Resources)
		if err != nil {
			return err
		}
	}
	return nil
}
