# Bulwark Vault

The primary motivation behind Bulwark Vault is to build a data storage system where the provider of the storage doesn't have access to the data and minimal access to user metadata. It does this by anonymizing most access using a system of login and access keys. The end result is that the provider cannot track users through the data directly provided to it (tracking IPs is a different matter, though).

This explanation was initially written as a blog post, but since that blog is no longer operational I figured this would be a good place to retain the explanations.

## Data Storage

![Data Storage Diagram](/images/diagram-1.png)

The basic foundation of the Vault is the actual data storage. We need our storage to be both private and anonymous, with no direct ties to our users. One easy step in this direction is to encrypt the data with encryption keys that only the client has access to: we can encrypt user data on the client using [AES-GCM](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) keys with 256 bits, which is fairly industry standard and securely encrypts the user data.

The more difficult question becomes, how do clients find and access this data without some sort of privacy-reducing authentication mechanism? To solve this, every encrypted blob is accessed using a uniformly distributed 256-bit access key that allows both read and write access to the data on the backend. Note that this access key does not need to be completely random; it can be derived from other data through cryptographic means, or can even be a SHA-256 hash of the data itself.

This data storage scheme, where data is encrypted with random encryption keys and accessed with random access keys, achieves our goals of privacy and anonymity. User data cannot be read by BulwarkID because it is encrypted, and it cannot be directly tracked because all accesses use individual, opaque access keys instead of authenticated requests.

It is important to note that while access keys and encrypted data will not expose any user information, certain metadata such as IP addresses and other fingerprinting mechanisms can work to break down user privacy by tracking users across requests. However, fixing such issues is difficult at this level and would be better solved by something like Tor to anonymize traffic.

## Access Keys and Encryption Keys

![Access Key Diagram](/images/diagram-2.png)

While we can now store and retrieve private data using random access and encryption keys, this isn’t a very convenient system of access. Every blob would need to be tracked and keys would need to be managed across devices. This is not an impossible problem to solve, but it would be useful to have a standard way of tracking these keys.

One solution is just to use a simple index. Whenever the client needs to store a new piece of data, generate a new access key and a new encryption key and add it to the index (perhaps a JSON blob). We can key this index with something useful, such as the full path of the data in a virtual filesystem. For example, when a client wants to access `/user/profile`, they can look up the access and encryption key for `/user/profile` in their index and then retrieve it from the backend.

This solution becomes even more useful if we store the index on the backend as well. We can generate a master access key and a master encryption key and then use that to store the JSON blob on the backend, which we can then retrieve to access any other data we have stored. This turns the multitude of access and encryption keys we need to track into a single pair.

## Master Keys and Key Derivation

![Key Derivation Diagram](/images/diagram-3.png)

The index solution works, but it has technical tradeoffs. Every vault access requires having an up to date version of the index, requiring some sort of synchronization mechanism between clients. As more and more keys are added, and as we potentially expand to more use cases such as multi-user data, this synchronization problem becomes more complex. It would be useful to be able to take a single piece of data, such as some root or master key, and deterministically turn it into multiple keys.

This is what we can do with key derivation. [HKDF](https://en.wikipedia.org/wiki/HKDF) is a cryptographic function that takes in some input material (plus a salt) and derives new, cryptographically secure key material out of it. The new material is still limited by the initial entropy of the input material, but this method does allow us to stretch some high entropy material into multiple usable keys.

Lets say we start with some high entropy material, such as a cryptographically random master secret of 256 bits. We can combine this secret with a file path and use HKDF to generate access keys and encryption keys. In the Vault, we use HKDF on input data of the form `object_access_key:<master secret>:<path>` in order to derive the access key and use input data of the form `object_encryption_key:<master secret>:<path>` to derive the encryption key. In this instance, path is a unique filepath that identifies the blob for that master secret.

In practice, this comes with its own limitations. The primary limitation is that there is one single, deterministic access and encryption key for each path and master secret (this is coincidentally also the main benefit of this system). This means that if we wanted to rotate keys, we would need to change the access key and encryption key for every piece of data in our system.

The Vault uses a combination of both. Stored indexes can be used for most objects, in order to allow for easy rotation of the keys, and key derivation will be used in a limited capacity where storing an index is not feasible. Namely, key derivation is used in order to access the index itself; the client can derive the index keys using a single 256-bit master secret and the path `/directory`.

## Logging In with Username and Password

![Login Diagram](/images/diagram-4.png)

While it is nice to be able to retrieve all user data from the system with just a single master key, this isn’t how most users want to be able to log in. This fact can be covered up somewhat by just storing master keys locally on the device, but that doesn’t help with the issue of syncing across devices. What most users want is to be able to log into the system with just a normal username and password.

There is a way to turn a username and password into a master key. PBKDF2 allows us to turn low-entropy material into uniformly distributed, though still low entropy, key material. It also makes is relatively difficult to reverse, since it is expensive to compute and thus expensive to brute force. PBKDF2 is commonly used when you need deterministic encryption based on low-entropy material, such as encrypting a vault in a password manager, so it is exactly what we need in our case.

In the Vault, we use [PBKDF2](https://en.wikipedia.org/wiki/PBKDF2) to derive a “login master secret” by inputting material of the form `login_secret:<username>:<password>` and deriving a 256-bit secret. We then use this login secret to derive access and encryption keys using HKDF, as above, to retrieve a high-entropy master secret stored on the backend.

Why retrieve a separate master secret from the backend and not just use the login secret as the master secret? For one, we would like to allow users to change their username and password without having to re-encrypt everything with the new login secret; it becomes easy just to change the login secret and store the old master secret at the new location.

Secondly, we want to contain the use of low entropy data as much as possible. PBKDF2 can make it difficult to reverse but doesn’t actually add entropy to the input material. By exchanging the low entropy login secret with a high entropy master secret on the backend, we can be more reasonably assured that the security guarantees of the rest of the system are maintained.
