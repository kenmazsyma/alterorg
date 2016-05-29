contract User {

	address		account=msg.sender;
	bytes[]		ipfsNodes;
	string		name;
	// TODO:implement oparation history

	function User(bytes node, string n) {
		name = n;
		ipfsNodes.push(node);
	}

	function isExistNode(bytes node) returns(bool) {
		for ( uint i=0; i<ipfsNodes.length; i++ ) {
			if (bytesEqual(ipfsNodes[i], node)) {
				return true;
			}
		}
		return false;
	}

	function changeName(string n) {
		name = n;
	}

	function appendIpfsNode(bytes node) {
		ipfsNodes.push(node);
	}
	
	// TODO:move to library
	function bytesEqual(bytes storage _a, bytes memory _b) internal returns (bool) {
		bytes storage a = bytes(_a);
		bytes memory b = bytes(_b);
		if (a.length != b.length)
			return false;
		// @todo unroll this loop
		for (uint i = 0; i < a.length; i ++)
			if (a[i] != b[i])
				return false;
		return true;
	}
}
