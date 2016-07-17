
function User(address) {
	this.address = address;
	this.regId = -1;
}

User.prototype.term = function() {
	if ( this.regId==-1) {
		window.clearInterval(this.regId);
		this.regId = -1;
	}
}

User.prototype.getInfo = function() {
	var _this = this;
	var get = function() {
			if (!_this.address) {
				_this.address = AlterOrg.checkMyAddress();
			}
			rpccall('User.GetInfo', [_this.address], function(res) {
				$('#name').val(res.result.name);
				$('#adrsusr').html(_this.address);
				$('#adrseth').html(res.result.adrs4eth);
				$('#adrsipfs').html(res.result.adrs4ipfs);
			},
			function(res, stat, err) {
				alert(err.message);
			});
		};
	var state = AlterOrg.checkMyAddress();
	if (state==AlterOrg.CODE_MYADDRESS_ERRORED) return;
	if (state==AlterOrg.CODE_MYADDRESS_SUCCESS) {
		AlterOrg.appendUserQue(get)
	} else {
		this.address = state;
		get();
	}
}

User.prototype.reg = function(name) {
	var _this = this;
	rpccall('User.Reg', [{node:'dummy', name:name}], function(res) {
		_this.regId = window.setInterval(function() {
			rpccall('User.CheckReg', [res.result], function(res2) {
				if (res2.result!="") {
					_this.address = res2.result;
					AlterOrg.updateMyAddress(_this.address);
					window.clearInterval(_this.regId);
					_this.getInfo();
				}
			},
			function(res, stat, err) {
				alert(err.message)
				window.clearInterval(_this.regId);
			})
		}, 1000);
	},
	function(res, stat, err) {
		alert(err.message)
	});
}


User.prototype.draw = function() {
	var _this = this;
	var fm = ''
	 + '<h1>UserInfo</h1>'
	 + '<div class="form-group">'
	 + '<label>name:</label><br>'
	 + '<input type="text" id="name" class="form-control" side="30">'
	 + '<label>address for User:</label><br>'
	 + '<div id="adrsusr">'
	 + '</div><br>'
	 + '<label>address for Ethereum:</label><br>'
	 + '<div id="adrseth">'
	 + '</div><br>'
	 + '<label>address for IPFS:</label><br>'
	 + '<div id="adrsipfs">'
	 + '</div><br>'
	 + '</div>'
	 + '<div class="form-group">'
	 + '<ul class="list-inline">'
	 + '<li><button class="btn btn-primary" id="btn_reg">Regist</button></li>'
	 + '<li><button class="btn btn-primary" id="btn_upd">Update</button></li>'
	 + '</ul>'
	 + '</div>'
	 + '</ul>';

	$('#main').html(fm);
	$('#btn_reg').click(function() {
		_this.reg($('#name').val());
	});
	$('#btn_upd').click(function() {
		_this.update($('#name').val());
	});
	this.getInfo();
}

