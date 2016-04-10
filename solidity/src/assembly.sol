contract Assembly {

	struct Proposal {
		bytes		docHash;
		bytes		discussHash;
		address		arbiter;
	}

	Proposal[]		proposal;	// history data for proposal address

	// constructor
	function Assembly(bytes hop, bytes hod) {
		revisionProposal(hop, hod);
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
