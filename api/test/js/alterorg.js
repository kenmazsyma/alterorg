var RpcSv = '';

function initSv(url){
	RpcSv = url;
}

var id = 1;
function rpccall(method, param, cb, cberr, failcb) {
	var data = {
		'jsonrpc':'2.0',
		'id':id++,
		'method':method,
		'params':param
	};
	if (failcb===failcb||type(failcb)!=='function') {
		failcb = function(jqXHR, textStatus, errorThrown) {};
	}
	$.ajax({
		type: "POST",
		url: RpcSv,
		data: JSON.stringify(data),
		dataType: 'json',
		contentType: "application/json",
		success: cb,
		error : cberr
	}).fail(failcb);
}

Vw = {
	breadCrumb : function(l) {
		var ret = '<ol class="breadcrumb">';
		for ( var i in l ) {
			if (i==l.length-1) {
				ret += '<li class="active">'
				+ l[i]['name']
				+ '</li>';
			} else {
				ret += '<li>'
				+ '<a href="#" onclick="'
				+ l[i].link
				+ '">'
				+ l[i].name
				+ '</a></li>';
			}
		}
		ret += '</ol>';
		$('#bread').html(ret);
	},

	home : function() {
		this.breadCrumb([{name:'Home'}]);
		$('#main').html('');
	},

	userinfo : function(adrs) {
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
		 + '</ul>'
		 + '</div>'
		 + '</ul>';

		$('#main').html(fm);
		$('#btn_reg').click(function() {
			user.reg($('#name').val());
		});
		if (adrs===undefined) {
			user.getMyData();
		} else {
			user.getInfo(adrs);
		}
	},
	board : function(address) {
		var fm = ''
		 + '<h1>Board Test</h1>'
		 + '<div class="form-group">'
		 + '<label>list</label>'
		 + '<div id="list">'
		 + '</div>'
		 + '</div>'
		 + '<div class="form-group">'
		 + '<label>Content</label>'
		 + '<textarea id="content" class="form-control" rows=10 side="30"></textarea>'
		 + '<ul class="list-inline">'
		 + '<li><button class="btn btn-primary" id="btn_wri">Write</button></li>'
		 + '</ul>'
		 + '</div>'
		 + '</ul>';
		$('#main').html(fm);
		$('#btn_wri').click(function() {
			Board.write($('#content').val());
		});
		this.breadCrumb([{name:'Home', link:'Vw.home()'}, {name:'Assembly detail'}, {name:'Board'}]);
		Board.init(address);
	}
};


function User() {
	this.orgLst = {};
}

User.prototype.getAssemblyList = function() {
	var _this = this;
	rpccall(
		'Alterorg.QueryAssemblyLst',
		[],
		function(res){
			_this.orgLst = {};
			for (var i in res.result ) {
				_this.orgLst[res.result[i]] = {}
			}
			_this.draw();
		},
		function(req, stat,err) {
			outputLog(err.message);
		}
	);
}

User.prototype.getInfo = function() {
	var _this = this;
	rpccall('User.GetInfo', [this.contAdrs], function(res) {
		$('#name').val(res.result.name);
		$('#adrsusr').html(_this.contAdrs);
		$('#adrseth').html(res.result.adrs4eth);
		$('#adrsipfs').html(res.result.adrs4ipfs);
	},
	function(res, stat, err) {
		alert(err.message);
	});
}

User.prototype.reg = function(name) {
	var _this = this;
	rpccall('User.Reg', [{node:'dummy', name:name}], function(res) {
		var intervalId;
		intervalId = window.setInterval(function() {
			rpccall('User.CheckReg', [res.result], function(res2) {
				if (res2.result!="") {
					_this.contAdrs = res2.result;
					window.clearInterval(intervalId);
					_this.getInfo();
				}
			},
			function(res, stat, err) {
				alert(err.message)
				window.clearInterval(intervalId);
			})
		}, 1000);
		
	},
	function(res, stat, err) {
		alert(err.message)
	});
}


User.prototype.getMyData = function() {
	this.contAdrs = '';
	var _this = this;
	rpccall('Alterorg.GetEthAddress', [], function(res) {
		if (!res.result||res.result==='') {
			alter("can't edit user info because of not started ethererum");
			return;
		} else {
			rpccall('User.GetMappedUser', [res.result], function(res2) {
				if (res2.result&&res2.result!=='') {
					_this.contAdrs=res2.result;
					_this.getInfo();
				}
			},
			function(res, stat, err) {
				alert(err.message)
			});
		}
	},
	function(res, stat, err) {
		alert(err.message)
	});
}

User.prototype.draw = function() {
	$('#orglist').html('');
	var code = '';
	for ( var i in this.orgLst ) {
		code += '<li role="presentation">'
			  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.Assemlby.detail('
			  + "'" + i + "'" 
			  + ')">'
			  + i
			  +'</a></li>';
	}
	code += '<li role="presentation">'
		  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.userinfo()">'
		  + 'Edit Userinfo'
		  +'</a></li>';
	code += '<li role="presentation">'
		  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.Assembly.new()">'
		  + 'Create new Assembly'
		  +'</a></li>';
	$('#orglist').append(code);
}


Board = {
	address : '',
	last : 1,
	init : function(address) {
		this.address = address;
		rpccall('Alterorg.PrepareBoard', [Board.address], function(res) {
		});
		window.setTimeout(Board.draw,1000);
	},
	draw : function() {
		rpccall('Alterorg.ListBoard', [Board.address], function(res) {
			var ret = '';
			var idx = 0;
			for (var i=0; i<res.result.length; i++) {
				ret += '<div>' + res.result[i][1].replace(/[\r]*\n/g, '<br>') + '</div>';
				idx = parseInt(res.result[i][0]);
				Board.last = (Board.last>idx) ? Board.last : idx;
			}
			$('#list').html(ret);
		}
		)
	},
	write : function (txt) {
		rpccall('Alterorg.WriteToBoard', [[Board.address, txt, (Board.last+1).toString()]], function(res) {
			Board.draw();
		},
		function(res, stat, err) {
			Board.draw();
			alert(err.message);
		},
		function() {
			alert('failed');
		}
		)
	}
}

