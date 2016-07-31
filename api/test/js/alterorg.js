AlterOrg = {
	CODE_MYADDRESS_ERRORED : '0x',
	CODE_MYADDRESS_NOT_YET : '',
	itrid  : -1,
	accnt : null,
	coinbase : '',
	userAdrs : this.CODE_MYADDRESS_NOT_YET,
	orgLst : {},
	usrQue : [],
	connect : function() {
		var _this = this;
		openModal('connecting to Alterorg server...');
		this.itrid = window.setInterval(
			function() {
				var term = function() {
					window.clearInterval(_this.itrid);
					_this.itrid = -1;
					closeModal();
				};
				rpccall(
					'Alterorg.GetEthStatus',
					[],
					function(res){
						if (res.result!=='RUN') {
							// wait for initialization
						} else {
							term();
							_this.receiveAssemblyList();
							_this.receiveMyAddress();
						}
					},
					function(req, stat,err) {
						term();
						if (stat=='error') {
							tempModal('Failed to connect AlterOrg server.');
						} else {
							tempModal(err.message);
						}
					}
				);
			}
		, 500);
	},

	receiveAssemblyList : function() {
		var _this = this;
		rpccall(
			'Alterorg.QueryAssemblyLst',
			[],
			function(res){
				_this.orgLst = {};
				for (var i in res.result ) {
					_this.orgLst[res.result[i]] = {}
				}
				_this.drawMenu();
			},
			function(req, stat,err) {
				outputLog(err.message);
			}
		);
	},

	receiveMyAddress : function() {
		var _this = this;
		rpccall('Alterorg.GetEthAddress', [], function(res) {
			if (!res.result||res.result==='') {
				alart("can't edit user info because of not started ethererum");
				return;
			} else {
				_this.coinbase = res.result;
				rpccall('User.GetMappedUser', [res.result], function(res2) {
					if (res2.result&&res2.result!=='') {
						_this.userAdrs=res2.result;
						for ( var i in _this.usrQue ) {
							i.f(i.o);
						}
					}
				},
				function(res, stat, err) {
					alert(err.message);
					_this.coinbase = AlterOrg.CODE_MYADDRESS_ERRORED;
				});
			}
		},
		function(res, stat, err) {
			alert(err.message);
			_this.coinbase = AlterOrg.CODE_MYADDRESS_ERRORED;
		});
		
	},

	drawMenu : function() {
		$('#orglist').html('');
		var code = '';
		for ( var i in this.orgLst ) {
			code += '<li role="presentation">'
				  + '<a role="menuitem" tabindex="-1" href="#" onclick="$M.assembly('
				  + "'" + i + "'" 
				  + ')">'
				  + i
				  +'</a></li>';
		}
		code += '<li role="presentation">'
			  + '<a role="menuitem" tabindex="-1" href="#" onclick="$M.user()">'
			  + 'Edit Userinfo'
			  +'</a></li>';
		code += '<li role="presentation">'
			  + '<a role="menuitem" tabindex="-1" href="#" onclick="$M.assembly(\'\', 1)">'
			  + 'Create new Assembly'
			  +'</a></li>';
		code += '<li role="presentation">'
			  + '<a role="menuitem" tabindex="-1" href="#" onclick="$M.assembly(\'\', 2)">'
			  + 'Join to Assembly'
			  +'</a></li>';
		$('#orglist').append(code);
	
	},

	draw : function() {
		$M.breadCrumb([{name:'Home'}]);
		$('#main').html('');
	},

	appendUserQue : function (f,o) {
		this.usrQue.push({f:f,o:o});
	},

	checkMyAddress : function() {
		return this.userAdrs;
	},

	updateMyAddress : function(address) {
		this.userAdrs = address;
	},

	term : function() {
	}

};

