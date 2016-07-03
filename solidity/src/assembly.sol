contract Assembly {

	struct Version {
		int idxName;
		int idxProposal;
		int[] idxTask; // -1 means "not using"
		int idxIgnoreList;
		address arbiter;
	}

	struct Proposal {
		bytes		docHash;
		bytes		discussHash;
		address		arbiter;
	}

	struct Participant {
		address		person;
		int			nofToken;
	}

	struct Task {
		string name;
		string desc;
		address pic;
		address[] reviewer;
		int		reward;
		uint	status;
		string	derivatives;
	}

	Version[]	 version;
	Proposal[]		proposal;	// history data for proposal address
	Task[][/*ver*/]	tasklist;
	address[/*ver*/][]	ignorelist;

	string[/*ver*/] name;
	string empty = "";

	// out of coverage for history
	Participant[]	participant;

	event onCreated(address adrs);
	// constructor
	function Assembly(string n) {
		name.push(n);
		version.length++;
		Version ver = version[version.length-1];
		ver.idxName = int(name.length-1);
		ver.idxProposal = -1;
		ver.idxIgnoreList = -1;
		ver.arbiter = msg.sender;
		participant.push(Participant({
			person:msg.sender,
			nofToken:0
		}));
		onCreated(this);
	}

	function getBasicInfo() returns (string, string, address, uint) {
		string nret =empty;
		string pret =empty;
		address vret = address(0x0);
		if (name.length>0) {
			nret = name[name.length-1];
		}
		if (proposal.length>0) {
			pret = string(proposal[proposal.length-1].docHash);
		}
		if (version.length>0) {
			vret = version[version.length-1].arbiter;
		}
		return (nret, pret, vret, version.length);
	}

	event onAddedPerson(address[] adrs);
	function addPerson(address[] adrs) {
		for (uint i=0; i<adrs.length; i++) {
			participant.push(Participant({
				person:adrs[i],
				nofToken:0
			}));
		}
		onAddedPerson(adrs);
	}

	function getName() returns (string) {
		return name[name.length-1];
	}
	
	// functions for revisioning proposal

	event onRevisionedProposal(address adrs, uint version);

	function revisionProposal(bytes hop, bytes hod) {
		proposal.push(Proposal({
							docHash:hop, 
							discussHash:hod,
							arbiter:msg.sender
					}));
		onRevisionedProposal(this, proposal.length-1);
	}

	// functions for refering proposal

	function getProposal() returns(bytes, bytes, address) {
		Proposal p = proposal[proposal.length-1];
		return (p.docHash, p.discussHash, p.arbiter);
	}

	function getProposalHistory(uint ver) returns(bytes, bytes, address){
		Proposal p = proposal[ver];
		return (p.docHash, p.discussHash, p.arbiter);
	}
}
