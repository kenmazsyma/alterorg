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

	newAssembly : function() {
		var fm = ''
		 + '<h1>Create new assembly</h1>'
		 + '<div class="form-group">'
		 + '<label>name</label>'
		 + '<input type="text" id="oname" class="form-control" side="30">'
		 + '</div>'
		 + '<input id="sel_fl" type="file" style="display:none">'
		 + '<label>proposal</label>'
		 + '<div class="input-group">'
		 + '<input type="text" id="path" class="form-control" placeholder="select file...">'
		 + '<span class="input-group-btn">'
		 + '<button type="button" class="btn btn-default" id="btn_fl">Browse</button>'
		 + '</span>'
		 + '</div>'
		 + '<br>'
		 + '<ul class="list-inline">'
		 + '<li><button class="btn btn-primary" id="btn_cre">Create</button></li>'
		 + '<li><button class="btn btn-default" id="btn_can">Cancel</button></li>'
		 + '</ul>';
		$('#main').html(fm);
		$('#sel_fl').change(function() {
			$('#path').val(this.files[0].name);
		});
		$('#btn_fl').click(function(){
			$('#sel_fl').click();
		});
		$('#btn_cre').click(function() {
			var asm = new Assembly()
			asm.new($('#oname').val())
		});
		$('#btn_can').click(function() {
			Vw.home();
		});
		this.breadCrumb([{name:'Home', link:'Vw.home()'}, {name:'Create new assembly'}]);
	},
	detailAssembly : function(address) {
		var _this = this;
		var fm = ''
		 + '<h1>Assembly Detail</h1>'
		 + '<div class="form-group">'
		 + '<label>name:</label>'
		 + '<span id="oname"></span>&nbsp;(&nbsp;'
		 + address
		 + '&nbsp;)&nbsp;<br>'
		 + '<label>proposal:</label>'
		 + '<a href="#" id="proposal"></a><br>'
		 + '<label>last arbiter:</label>'
		 + '<a href="#" id="arbiter"></a><br>'
		 + '<label>version:</label>'
		 + '<a href="#" id="version"></a>'
		 + '</div>'
		 + '<br>'
		 + '<ul class="list-inline">'
		 + '<li><button class="btn btn-primary" id="btn_board">Board</button></li>'
		 + '<li><button class="btn btn-primary" id="btn_close">Close</button></li>'
		 + '</ul>';
		$('#main').html(fm);
		$('#btn_board').click(function() {
			Vw.board(address);
		});
		$('#btn_close').click(function() {
		});
		this.breadCrumb([{name:'Home', link:'Vw.home()'}, {name:'Assembly detail'}]);
		rpccall('Assembly.GetBasicInfo', [address], function(res) {
			$('#oname').html(res.result.name);
			$('#arbiter').html(res.result.arbiter);
			$('#version').html(res.result.version);
		},
		function(res, stat, err) {
			alert(err.message)
		});
	
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


function Assembly(address) {
	this.address = address;
}

Assembly.prototype.new = function(name) {
	var _this = this;
	rpccall('Assembly.Create', [name], function(res) {
		var intervalId;
		intervalId = window.setInterval(function() {
			rpccall('Assembly.CheckMine', [res.result], function(res2) {
				if (res2.result!="") {
					_this.address = res2.result;
					rpccall('Alterorg.AppendAssembly', [res2.result], function(res3) {
						user.getAssemblyList();
						window.clearInterval(intervalId);
					});
				}
			},
			function(res, stat, err) {
				alert(err.message)
				window.clearInterval(intervalId);
			})
		}, 1000);
	}
	)
}

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

User.prototype.draw = function() {
	$('#orglist').html('');
	var code = '';
	for ( var i in this.orgLst ) {
		code += '<li role="presentation">'
			  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.detailAssembly('
			  + "'" + i + "'" 
			  + ')">'
			  + i
			  +'</a></li>';
	}
	code += '<li role="presentation">'
		  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.board()">'
		  + 'Board test'
		  +'</a></li>';
	code += '<li role="presentation">'
		  + '<a role="menuitem" tabindex="-1" href="#" onclick="Vw.newAssembly()">'
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

