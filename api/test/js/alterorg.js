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
		 + '<input type="text" name="oname" class="form-control" side="30">'
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
			alert('create');
		});
		$('#btn_can').click(function() {
			Vw.home();
		});
		this.breadCrumb([{name:'Home', link:'Vw.home()'}, {name:'Create'}]);
	}
}