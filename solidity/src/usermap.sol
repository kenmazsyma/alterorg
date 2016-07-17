import "../src/user.sol";

contract UserMap {

	mapping(address=>User) usermap;
	address[] list;

	event onReg(address adrs, address cont, bool isNew);

	function reg(bytes node, string n) {
		User user = new User(msg.sender, node, n);
		if (usermap[msg.sender]==address(0x0)) {
			list.push(msg.sender);
			usermap[msg.sender] = user;
			onReg(msg.sender, address(user), true);
		} else {
			usermap[msg.sender] = user;
			onReg(msg.sender, address(user), false);
		}
	}

	function getAddresses() returns(address[]){
		return list;
	}

	function getUser(address adrs) constant returns(address){
		return address(usermap[adrs]);
	}

	function getName(address adrs) returns (address) {
		return User(adrs).getName();
	}

}

