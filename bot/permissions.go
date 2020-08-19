package heimdallr

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//CheckPermissions checks if Heimdallr has the permissions it needs
func CheckPermissions(s *discordgo.Session) {
	for i := 0; i < 10; i++ {
		adminLogChannelPermissions, err := s.State.UserChannelPermissions(s.State.User.ID, Config.AdminLogChannel)
		if err != nil {
			if i < 9 {
				time.Sleep(time.Second * 10)
				continue
			}
			log.Fatalf("%+v\n", errors.Wrap(err, "couldn't fetch permissions for admin log channel"))
		}
		welcomeChannelPermissions, err := s.State.UserChannelPermissions(s.State.User.ID, Config.WelcomeChannel)
		if err != nil {
			if i < 9 {
				time.Sleep(time.Second * 10)
				continue
			}
			log.Fatalf("%+v\n", errors.Wrap(err, "couldn't fetch permissions for welcome channel"))
		}

		if PermissionInPermissions(discordgo.PermissionAdministrator, adminLogChannelPermissions) {
			return
		}
		if !CanWrite(adminLogChannelPermissions) {
			log.Fatalf("%+v\n", errors.Wrap(err, "the bot lacks permissions to write in the admin log channel"))
		}

		if !PermissionInPermissions(discordgo.PermissionBanMembers, adminLogChannelPermissions) {
			LogMessage(s, "The bot lacks permissions to ban, the `ban` command will not work.")
		}

		if !PermissionInPermissions(discordgo.PermissionKickMembers, adminLogChannelPermissions) {
			LogMessage(s, "The bot lacks permissions to kick, the `kick` command will not work.")
		}

		if !PermissionInPermissions(discordgo.PermissionManageRoles, adminLogChannelPermissions) {
			LogMessage(s, "The bot lacks permissions to manage roles, the `role` and `approve` commands will not work.")
		}

		if !PermissionInPermissions(discordgo.PermissionCreateInstantInvite, welcomeChannelPermissions) {
			LogMessage(s, "The bot lacks permissions to create invites to the welcome channel, the `invite` command will not work.")
		}

		if !CanWrite(welcomeChannelPermissions) {
			LogMessage(s, "The bot lacks permissions to write in the welcome channel, the welcome and approval messages will not work.")
		}
		break
	}
}

//hasRole checks if a Member has the specified role.
func hasRole(m *discordgo.Member, r string) bool {
	for _, role := range m.Roles {
		if role == r {
			return true
		}
	}
	return false
}

//IsTrialModOrHigher returns whether the member has mod permissions or higher
func IsTrialModOrHigher(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.TrialModRole) || IsModOrHigher(member, guild)
}

//IsModOrHigher returns whether the member has mod permissions or higher
func IsModOrHigher(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.ModRole) || IsSuperModOrHigher(member, guild)
}

//IsSuperModOrHigher returns whether the member has supermod permissions or higher
func IsSuperModOrHigher(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.SuperModRole) || IsAdminOrHigher(member, guild)
}

//IsAdminOrHigher returns whether the member has admin permissions or higher
func IsAdminOrHigher(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.AdminRole) || IsOwner(member, guild)
}

//IsOwner returns whether the member is the guild's owner
func IsOwner(member *discordgo.Member, guild *discordgo.Guild) bool {
	return guild.OwnerID == member.User.ID
}

//IsHelper returns whether the member is a brothers or sisters helper
func IsHelper(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.BrothersHelperRole) || hasRole(member, Config.SistersHelperRole)
}

//IsCircleMember returns whether the member is in a circle or not
func IsCircleMember(member *discordgo.Member, guild *discordgo.Guild) bool {
	return hasRole(member, Config.OmerIbnAlKhattabRole) || hasRole(member, Config.AbuBakrAlSiddeeqRole) || hasRole(member, Config.AliBinAbiTaalibRole) || hasRole(member, Config.SistersCircleRole)
}

//IsVerified checks whether a user is verified
func IsVerified(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == Config.VerifiedMaleRole || role == Config.VerifiedFemaleRole {
			return true
		}
	}
	return false
}
func getRoleByID(roleID string, roles []*discordgo.Role) (*discordgo.Role, error) {
	for _, role := range roles {
		if roleID == role.ID {
			return role, nil
		}
	}
	return nil, errors.New("role not found")
}

func computeBasePermissions(member *discordgo.Member, guild *discordgo.Guild) int {
	if IsOwner(member, guild) {
		return discordgo.PermissionAll

	}
	roleEveryone, _ := getRoleByID(guild.ID, guild.Roles)
	permissions := roleEveryone.Permissions

	for _, roleID := range member.Roles {
		role, _ := getRoleByID(roleID, guild.Roles)
		permissions |= role.Permissions
	}

	if PermissionInPermissions(discordgo.PermissionAdministrator, permissions) {
		return discordgo.PermissionAll
	}

	return permissions
}

//PermissionInPermissions returns whether the permission number contains the given permission
func PermissionInPermissions(permission int, permissions int) bool {
	return permission&permissions == permission
}

//CanRead returns whether the permission number contains read permission
func CanRead(permissions int) bool {
	return PermissionInPermissions(discordgo.PermissionAdministrator, permissions) ||
		PermissionInPermissions(discordgo.PermissionReadMessages|discordgo.PermissionReadMessageHistory, permissions)
}

//CanWrite returns whether the permission number contains write permission
func CanWrite(permissions int) bool {
	return PermissionInPermissions(discordgo.PermissionAdministrator, permissions) ||
		PermissionInPermissions(discordgo.PermissionReadMessages|discordgo.PermissionSendMessages, permissions)
}
