function Assembly(address) {
	this.address = address;
	this.intervalId = -1;
}

Assembly.prototype.term = function() {
	if (this.intervalId!=-1) window.clearInterval(this.intervalId);
	this.intervalId = -1;
}

Assembly.prototype.new = function(name) {
	var _this = this;
	var term = function() {
		window.clearInterval(_this.intervalId);
		_this.intervalId = -1;
	};
	rpccall('Assembly.Create', [name], function(res) {
		_this.intervalId = window.setInterval(function() {
			rpccall('Assembly.CheckMine', [res.result], function(res2) {
				if (res2.result!="") {
					_this.address = res2.result;
					rpccall('Alterorg.AppendAssembly', [res2.result], function(res3) {
						AlterOrg.receiveAssemblyList();
						term();
					});
				}
			},
			function(res, stat, err) {
				alert(err.message);
				term();
			})
		}, 1000);
	}
	)
}

Assembly.prototype.draw = function() {
	var _this = this;
	var fm = [];
	if (!this.address) {
		fm.push('<h1>Create new Assembly</h1>');
	} else {
		fm.push('<h1>Assembly detail</h1>');
	}
	if (this.address) {
		fm.push(
		   '<div class="form-group">'
		 + '<div class="form-group"><label>name</label><br><span id="oname"></span></div>'
		 + '<div class="form-group"><label>proposal</label><br>'
		 + '<a href="#" id="proposal">Download</a></div>'
		 + '<div class="form-group"><label>arbiter</label><br><span id="arbiter"></span></div>'
		 + '<div class="form-group"><label>version</label><br><span id="version"></span></div>'
		 + '<ul class="list-inline">'
		 + '<li><button class="btn btn-default" id="btn_board">Board</button></li>'
		 + '<li><button class="btn btn-default" id="btn_edit">Edit</button></li>'
		 + '<li><button class="btn btn-primary" id="btn_close">Close</button></li>'
		 + '</ul>'
		);
		$('#main').html(fm.join(''));
		$('#btn_board').click(function() {
			$M.board(_this.address);
		});
		$('#btn_close').click(function() {
			$M.home();
		});
		$('#proposal').click(function() {
			alert('test');
		});
		$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Assembly detail'}]);
		rpccall('Assembly.GetBasicInfo', [this.address], function(res) {
			$('#oname').html(res.result.name);
			$('#arbiter').html(res.result.arbiter);
			$('#version').html(res.result.version);
		},
		function(res, stat, err) {
			alert(err.message)
		});
		rpccall('Assembly.GetParticipants', [this.address], function(res) {
			if (res.result&&res.result.persons) {
				var html = '';
				for ( var i=0; i<res.result.persons.length; i++ ) {
					html += '<a href="#">' + res.result.persons[i] + '</a><br>'
				}
				$('#participants').html(html);
			}
		},
		function(res, stat, err) {
			alert(err.message)
		});
	} else {
		fm.push ( 
		   '<div class="form-group">'
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
		 + '</ul>'
		);
		$('#main').html(fm.join(''));
		$('#sel_fl').change(function() {
			$('#path').val(this.files[0].name);
		});
		$('#btn_fl').click(function(){
			$('#sel_fl').click();
		});
		$('#btn_cre').click(function() {
			_this.new($('#oname').val())
		});
		$('#btn_can').click(function() {
			$M.home();
		});
		$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Create new assembly'}]);
	}
}
