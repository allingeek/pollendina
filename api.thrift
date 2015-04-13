# Pollendina thrift interface definition
# Jeff Nickoloff (jeff@allingeek.com)

namespace go pollendina
namespace cpp pollendina
namespace java pollendina
namespace perl pollendina

struct ServiceLocation {
	1: string sid,
	2: string hostname,
	3: string cid,
}

exception UnauthorizedCertificate {
	1: string subject,
}

exception NoSuchServiceLocation {
	1: string subject,
}

service Pollendina {
	// Health checks can call ping to test service availability
	void ping(),

	// Authorize tells pollendina to sign requests with the specified matching service location attributes
	void authorize(1:ServiceLocation location, 2:bool override),

	// Sign a CSR, or throw an exception if the CSR is not authorized.
	string sign(1:string csr) throws (1:UnauthorizedCertificate cert),

	// Get certificate for a specific service location
	string lookup(1:ServiceLocation location) throws (1:NoSuchServiceLocation dne)

	// Maybe include this for bootstrapping new machines?
	// string getCACert(),
}
