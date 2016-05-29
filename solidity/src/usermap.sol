import "../src/user.sol";

contract UserMap {

	mapping(address=>User) usermap;
	address[] list;

	event onReg(address adrs, address cont, bool isNew);

	function reg(bytes node, string n) {
		User user = new User(node, n);
		usermap[msg.sender] = user;
		if (usermap[msg.sender]==address(0x0)) {
			list.push(msg.sender);
			onReg(msg.sender, user, true);
		} else {
			onReg(msg.sender, user, false);
		}
	}

	function getAddresses() returns(address[]){
		return list;
	}


	function getUser(address adrs) returns(address){
		return usermap[adrs];
	}

}

