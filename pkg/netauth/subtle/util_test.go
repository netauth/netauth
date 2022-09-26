package subtle

import (
	"testing"
)

func TestCompareSSHKeys(t *testing.T) {
	keys := []string{
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCfTQ+DfECHAwzpLI9kgrKq3ARTOArbqGJC4U2gwkmP5yXmWc1KnpVQUNhGmoz6U7bRoTYL7vWIt0V7VgOq+QZ0bNb4p0YXmT5wcF6C1GBfioO3/0kjWSbIuIXzcg6mUXsJF+vARy7MzbXYPJo/ZAT7NwKmslmNEiIJyNfbrey2hVXnDuKYxlek2z0L9F+E4hythwXWvOONwpUp044Fm1vjdnaRkSwpZIddkPbHCjvBPjrI9cNdmNLCuEZKflbSfsmfnAhPHvJ+rgtvuWX9A89ieKcYzmg23dGQyOjk/c4iPnWhx76ZV8bRXvHGuLAMUKI1b031ZIfZIHW9Csa7m/lWuX3pqsm/J+1dRphnnfr0oocC1MAatyIhuiOYhg+OVNkQQacNSkBqQP8eR7u6shgRu88MLEHRrzPvlZVPHZuWoOWU7jZMfBH8Nr9517VgX2hh5hyStOEQtowjfItw4JM4naiMK4exxxDN7QuVmzTamGqxg3fiz+7jdIAFsFmwMX0=",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBzarZCQgcbnvGxpzjItruVue5R5a4wqP6dWDxCnEUBG blah blah blah",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC+AiVs2GWw62oszzYHqwgirBqleT2J3/J9wwn4uOjAW0IfEvPnZ2eIdIuRrYOscGd87HJmH0Px/+C+8WRhless2UwqT6W5RuqDkBPnz7yOyk8bpsSWQAKOe6BpDQAcreOGi7ocR7AV0J9mGb0ZN221lTUbPNImeoNblcKBzk21GQ==",
		"",
	}

	cases := []struct {
		key     string
		wantErr error
	}{
		{
			// in keys
			key:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCfTQ+DfECHAwzpLI9kgrKq3ARTOArbqGJC4U2gwkmP5yXmWc1KnpVQUNhGmoz6U7bRoTYL7vWIt0V7VgOq+QZ0bNb4p0YXmT5wcF6C1GBfioO3/0kjWSbIuIXzcg6mUXsJF+vARy7MzbXYPJo/ZAT7NwKmslmNEiIJyNfbrey2hVXnDuKYxlek2z0L9F+E4hythwXWvOONwpUp044Fm1vjdnaRkSwpZIddkPbHCjvBPjrI9cNdmNLCuEZKflbSfsmfnAhPHvJ+rgtvuWX9A89ieKcYzmg23dGQyOjk/c4iPnWhx76ZV8bRXvHGuLAMUKI1b031ZIfZIHW9Csa7m/lWuX3pqsm/J+1dRphnnfr0oocC1MAatyIhuiOYhg+OVNkQQacNSkBqQP8eR7u6shgRu88MLEHRrzPvlZVPHZuWoOWU7jZMfBH8Nr9517VgX2hh5hyStOEQtowjfItw4JM4naiMK4exxxDN7QuVmzTamGqxg3fiz+7jdIAFsFmwMX0=",
			wantErr: nil,
		},
		{
			// in keys but has a comment
			key:     "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBzarZCQgcbnvGxpzjItruVue5R5a4wqP6dWDxCnEUBG",
			wantErr: nil,
		},
		{
			// uncommented version in keys
			key:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC+AiVs2GWw62oszzYHqwgirBqleT2J3/J9wwn4uOjAW0IfEvPnZ2eIdIuRrYOscGd87HJmH0Px/+C+8WRhless2UwqT6W5RuqDkBPnz7yOyk8bpsSWQAKOe6BpDQAcreOGi7ocR7AV0J9mGb0ZN221lTUbPNImeoNblcKBzk21GQ== this is a comment",
			wantErr: nil,
		},
		{
			// not in keys
			key:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCvRbRsAZfw9F0FRmgQhxOnv97TIHleU3a6RzUVGTnekJox+hE8LRKNuFATxFaj5Zxufmpm308dYA6RtqCtBcf2swBXYf0y43RdiflMAsbmiAoNnQVV8W6WP8gBWWkS6ri9NBD0b8ezYKOF/w60oUalSgMYiE79pc2bx0DO2Ixw6w==",
			wantErr: ErrNoMatchingKeys,
		},
		{
			// empty key
			key:     "",
			wantErr: ErrNoMatchingKeys,
		},
	}

	for i, c := range cases {
		if err := CompareSSHKeys(keys, c.key); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
