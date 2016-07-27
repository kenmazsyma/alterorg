var RpcSv = '';

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
		url: RpcSv + '/rpc',
		data: JSON.stringify(data),
		dataType: 'json',
		contentType: "application/json",
		success: cb,
		error : cberr
	}).fail(failcb);
}

function openModal(message){
	$('#modalmes').html(message);
	$('#modal').modal({backdrop: 'static', keyboard: false});
}

function closeModal() {
	$("#modal").modal('hide');
}

function tempModal(message) {
	openModal(message);
	window.setTimeout(closeModal, 3000);
}

function outputLog(msg) {
	console.log(msg);
}

$M = {
	cur : null,
	init : function(url) {
		RpcSv = url;
		AlterOrg.connect();
		this.cur = AlterOrg;
	},
	draw : function() {
		this.cur.draw();
	},
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
	changePage : function(o) {
		this.cur.term();
		this.cur = o;
	},
	assembly : function(address) {
		this.changePage(new Assembly(address));
		this.draw();
	},
	user : function(address) {
		this.changePage(new User(address));
		this.draw();
	},
	home : function() {
		this.changePage(AlterOrg);
		this.draw();
	},
	board : function(address) {
		this.changePage(new Board(address));
		this.draw();
	}
};


function QElm(n, p, f1, f2) {
	this.n = n;
	this.p = p;
	this.f1 = f1;
	this.f2 = f2;
	this.use = true;
	this.drawf = null;
	this.evt = null;
}

QElm.prototype.setDrawCB = function(f1, f2) {
	this.drawf = f1;
	this.evt = f2;
}

QElm.prototype.end = function() {
	this.use = false;
}

function QSet(interval) {
	this.use = true;
	this.lst = [];
}

QSet.prototype.start = function(interval) {
	this.interval = interval;
	$Q.start(this);
}

QSet.prototype.end = function() {
	this.use = false;
}

QSet.prototype.append = function(elm) {
	this.lst.push(elm);
}

$Q = {
	id : -1,
	lst : [],
	start : function(set) {
		this.lst.push(set);
		var _id = -1;
		var proc = function(){
			var term = function() {
				$Q.progress();
			};
			if (!set.use) {
				window.clearInterval(_id);
				return term();
			}
			if (set.lst.length==0) return term();
			if (!set.lst[0].use) {
				set.lst.shift();
			}
			if (set.lst.length==0) return term();
			var elm = set.lst[0];
			if (typeof(elm.n)=='function') {
				elm.n(elm.p, elm.f1, elm.f2);
			} else {
				rpccall(elm.n, elm.p, elm.f1, elm.f2);
			}
			term();
		};
		proc();
		_id = window.setInterval(proc, set.interval);
	},
	progress : function() {
		var html = [];
		var evt = {};
		for ( var i=0; i<this.lst.length; i++ ) {
			if (!this.lst[i].use) {
				this.lst.splice(i--,1);
				continue;
			}
			for ( var j=0; j<this.lst[i].lst.length; j++ ) {
				var elm = this.lst[i].lst[j];
				if (elm.use&&elm.drawf) {
					var id = 'prog_' + i + '_' + j;
					html.push(elm.drawf(id));
					evt[id] = elm.evt;
				}
			}
		}
		$('#prog').html(html.join(''));
		for ( var i in evt) {
			$('#' + i).click(evt[i]);
		}
	}
};


