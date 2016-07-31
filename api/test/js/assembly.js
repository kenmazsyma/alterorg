function Assembly(address, pat) {
	this.address = address;
	this.intervalId = -1;
	if (pat===2) {
		this.mode = 'JOIN';
	} else if (pat===1) {
		this.mode = 'ADD';
	} else {
		this.mode = 'REF';
	}
	this.prophash = '';
	this.propname = '';
	this.qset = new QSet();
	this.qset.start(2000);
	this.pat = pat;
}

Assembly.prototype.term = function() {
	if (this.intervalId!=-1) window.clearInterval(this.intervalId);
	this.intervalId = -1;
	this.qset.end();
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

Assembly.prototype.update = function(name, path) {
	var _this = this;
	var file = $('#sel_fl').prop('files')[0];
	var reader = new FileReader();
	reader.onloadend = function () {
		rpccall('Assembly.RevisionProposal', [{address:_this.address, discussion:'', propname:file.name, propdata:reader.result}], function(res) {
			var elm = new QElm('Assembly.CheckRevisionProposal', [res.result], function(res2) {
				if (res2.result&&res2.result.address!="") {
					alert(res2.result.address);
					elm.end();
				}
			}, function(req, stat, err) {
			});
			_this.qset.append(elm);
		}, function(req, stat, err) {
		});
	}
	reader.readAsDataURL(file);
}

Assembly.prototype.download = function() {
	var childWindow = window.open('about:blank');
	rpccall('Alterorg.QuerySetting', [['download_dir']], function(res) {
		rpccall('Alterorg.GetFile', [{hash:this.prophash, path:res.result + '/' + this.propname}], function(res2) {
			childWindow.location.href = 'file://' + res.result;
			childWindow = null;
		},
		function(res2, stat, err) {
			childWindow.close();
			childWindow = null;
			alert(err.error);
		});
	},
	function(res, stat, err) {
		alert(err.error);
	});
}

Assembly.prototype.draw = function() {
	var _this = this;
	var fm = [];
	var fb = {
		'REF' : function() {
			fm.push(
			   '<h1>Assembly detail</h1>'
			 + '<div class="form-group">'
			 + '<div class="form-group"><label>name</label><br><span id="oname"></span></div>'
			 + '<div class="form-group"><label>proposal</label><br>'
			 + '<span id="propname"></span>&nbsp;(&nbsp;<a href="#" id="prophash">'
			 + '</a>&nbsp;)&nbsp;</div>'
			 + '<div class="form-group"><label>arbiter</label><br><span id="arbiter"></span></div>'
			 + '<div class="form-group"><label>version</label><br><span id="version"></span></div>'
			 + '<div class="form-group"><label>participant</label><br><span id="participants"></span></div>'
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
			$('#btn_edit').click(function() {
				_this.edit = true;
				$M.draw();
			});
			$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Assembly detail'}]);
			rpccall('Assembly.GetBasicInfo', [_this.address], function(res) {
				_this.prophash = res.result.prophash;
				_this.propname = res.result.propname;
				$('#oname').html(res.result.name);
				$('#propname').html(res.result.propname);
				$('#prophash').html(res.result.prophash);
				$('#arbiter').html(res.result.arbiter);
				$('#version').html(res.result.version);
				$('#prophash').click(function() {
					_this.download();
				});
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
		},
		'EDIT' : function() {
			var _this = this;
			// TODO:sanitize
			fm.push(
			   '<h1>Assembly edit</h1>'
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
			 + '<li><button class="btn btn-primary" id="btn_upd">Update</button></li>'
			 + '<li><button class="btn btn-default" id="btn_can">Cancel</button></li>'
			 + '</ul>'
			);
			var buf_oname = $('oname').val();
			$('#main').html(fm.join(''));
			$('#sel_fl').change(function() {
				$('#path').val(this.files[0].name);
			});
			$('#btn_fl').click(function(){
				$('#sel_fl').click();
			});
			$('#btn_upd').click(function() {
				_this.update($('#oname').val(), $('#path').val());
			});
			$('#btn_can').click(function() {
				_this.mode = 'EDIT';
				_this.draw()
			});
			$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Create new assembly'}]);
		},
		'ADD' : function() {
			fm.push ( 
			   '<h1>Create new Assembly</h1>'
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
		},
		'JOIN' : function() {
			var _this = this;
			fm.push ( 
			   '<h1>Join to Assembly</h1>'
			 + '<div class="form-group">'
			 + '<label>address</label>'
			 + '<input type="text" id="adrs" class="form-control" side="30">'
			 + '<button class="btn btn-default" id="btn_search">Search</button>'
			 + '</div><br>'
			 + '<div class="form-group">'
			 + '<div class="form-group"><label>name</label><br><span id="oname"></span></div>'
			 + '<div class="form-group"><label>proposal</label><br>'
			 + '<span id="propname"></span>&nbsp;(&nbsp;<a href="#" id="prophash">'
			 + '</a>&nbsp;)&nbsp;</div>'
			 + '<div class="form-group"><label>arbiter</label><br><span id="arbiter"></span></div>'
			 + '<div class="form-group"><label>version</label><br><span id="version"></span></div>'
			 + '<div class="form-group"><label>participant</label><br><span id="participants"></span></div>'
			 + '<ul class="list-inline">'
			 + '<li><button class="btn btn-default" id="btn_join">Join</button></li>'
			 + '<li><button class="btn btn-primary" id="btn_close">Close</button></li>'
			 + '</ul>'
			);
			$('#main').html(fm.join(''));
			$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Join to Assembly'}]);
			$('#btn_search').click(function(){
				rpccall('Assembly.GetBasicInfo', [$('#adrs').val()], function(res) {
					if (res.error) {
						tempModal(res.error);
					} else {
						_this.prophash = res.result.prophash;
						_this.propname = res.result.propname;
						$('#oname').html(res.result.name);
						$('#propname').html(res.result.propname);
						$('#prophash').html(res.result.prophash);
						$('#arbiter').html(res.result.arbiter);
						$('#version').html(res.result.version);
						$('#prophash').click(function() {
							_this.download();
						});
					}
				},
				function(res, stat, err) {
					alert(err.message)
				});
				rpccall('Assembly.GetParticipants', [$('#adrs').val()], function(res) {
					if (!res.error) {
						if (res.result&&res.result.persons) {
							var html = '';
							for ( var i=0; i<res.result.persons.length; i++ ) {
								html += '<a href="#">' + res.result.persons[i] + '</a><br>'
							}
							$('#participants').html(html);
						}
					}
				},
				function(res, stat, err) {
					alert(err.message)
				});
			});
			$('#btn_join').click(function(){
				var adrs = $('#adrs').val();
				rpccall('Assembly.Join', [adrs], function(res) {
					if (res.error) {
						tempModal(res.error);
					} else {
						tempModal('success to join Assembly:' + adrs);
						$M.home();
						AlterOrg.receiveAssemblyList();
					}
				},
				function(res, stat, err) {
					alert(err.message)
				});
			});
			$('#btn_close').click(function() {
				$M.home();
			});
		}
	};
	fb[this.mode]();
}
