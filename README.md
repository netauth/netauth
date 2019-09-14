NetAuth
=======
NetAuth is a network identity and authentication provider.  It allows
you to have one user account that is available to a lot of different
machines.

The ultimate goal is to have a small service which could live in a
small VM and provide fleet wide authentication and identity services
for a small fleet of machines.

What Does it Do?
----------------

If you're familiar with LDAP and Kerberos, you can skip down to the
next section, NetAuth is an implementation of the services that LDAP
and Kerberos can provide for a network, but with a much smaller scope
and certain assumptions.

NetAuth provides two key components: a limited directory of user
information, and a secrets store.  The directory provides the most
critical information about an entitiy such as the ID, numeric ID,
name, etc.  What NetAuth does NOT provide is a general purpose
directory.  That is something which is not really in scope for a small
authentication service and is implemented exceptionally well by the
LDAP standard.  If you want such a server you should really setup
LDAP, which you could either use for authentication (something it was
not designed for) or use it as a directory that just contains
information.  If you need to authenticate your access to LDAP it would
not be too difficult to back up LDAP into NetAuth, but this
functionality is left as an exercise to the reader.

For identity NetAuth provides a fairly standard password verification
system that is not unlike that used by a website login system.  The
user's password is sent via a secure channel to the NetAuth server
where it is validated against a hashed copy.  If the password checks
out, then NetAuth will return a success message to the calling client.
In failure cases NetAuth will return a message to the client
explaining the failure.

Why is this written in Go?
--------------------------

I like Go and it works well with protobuf without needing the host
operating system to have good support on its own.  Its not the most
ideal language for interfacing with PAM or nsswitch, but for writing
servers that work with gRPC its quite nice.

Why does this communicate using gRPC and not my favorite protocol?
------------------------------------------------------------------

I like the RPC paradigm, it works well for what I am trying to achieve
here and can work without any real thought about the transport.  gRPC
specifically is capable of working in an environment where the only
outbound connection allowed is HTTP, which is a core design goal of
this project.  While it is a binary protocol, the protobuf definition
is public and this will let you do all the normal things such as
debugging with wireshark (assuming you have the appropriate security
settings in place to observe HTTPS traffic).

How do I hook up other things to NetAuth?
-----------------------------------------

There are several systems available to plug in to NetAuth.  For Linux
hosts you can use
[pam_netauth](https://github.com/NetAuth/pam_netauth) and
[nsscache](https://github.com/NetAuth/nsscache).  If you want to pull
ssh keys, then you probably want
[NetKeys](https://github.com/NetAuth/NetKeys).

If you use Okta, you'll probably be interested in the
[Okta Plugin](https://github.com/NetAuth/plugin-okta) which can
automatically mirror your NetAuth entities and groups into Okta.

Other modules are coming, if you want to help out, reach out in
`#netauth-dev` on freenode.

Why wouldn't you use LDAP and Kerberos?  Why did you build this?
----------------------------------------------------------------

I managed a network that used LDAP and Kerberos for a number of years.
These are some incredible technologies and I quite enjoyed the feature
sets they provide.  The problem though is that LDAP is a slimmed down
version of a protocol so complex it could never be implemented (the
DAP), and Kerberos is a protocol that makes certain assumptions about
the state of the network and the services that are available.  Both
LDAP and Kerberos require tooling to interact with and as far and wide
as I have searched, I have found no good tooling that allows one to
interact with the two as a single identity management platform (likely
because they aren't).

While I am a die-hard FOSS supporter and contributor, Microsoft's
Active Directory is by my book the gold standard of tooling for
managing an authentication and identity provider.  It is intuitive,
all-in-one and most importantly, it makes the underlying LDAP and
Kerberos servers behave as a single virtual service.  I thought long
and hard about whether or not I wanted to just build a frontend to
LDAP and Kerberos to handle the managerial tasks, but I came to the
conclusion that if I was going to write software from the ground up, I
would like to just rebuild the entire stack as a slimmed down version
that would do exactly what I wanted.  If nothing else I will gain a
strong appreciation for the work done by the developers of Kerberos
and LDAP.

If I was still managing networks with thousands of users I would
probably still stand up LDAP and Kerberos.  If I needed a directory of
arbitrary information, LDAP is still my first stop.  If I needed to do
very novel and interesting crypto to secure the network fabric I would
still use Kerberos.  For simple authentication on my home network or
on the Open Source projects I'm involved in, however, these
technologies are overkill.
