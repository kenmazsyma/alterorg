function Board(address) {
	this.address = address;
	var _this = this;
	this.qset = new QSet();
	var elm1 = new QElm('Alterorg.PrepareBoard', [this.address], function(res) {
		if (res.error) {
			tempModal(res.error);
			return;
		}
		elm1.end();
	}, function(res, stat, err) {
		alert(err);
		_this.qset.end();
	});
	var elm2 = new QElm('Alterorg.ListBoard', [this.address], function(res) {
		var ret = '';
		var idx = 0;
		if (res.error) {
			tempModal(res.error);
			return;
		}
		for (var i=0; i<res.result.length; i++) {
			ret += '<div>' + res.result[i][1].replace(/[\r]*\n/g, '<br>') + '</div>';
			idx = parseInt(res.result[i][0]);
			_this.last = (_this.last>idx) ? _this.last : idx;
		}
		$('#list').html(ret);
		_this.setDisable(false);
		$Q.progress();
	});
	elm1.setDrawCB(function(id) {
		return '<a href="#" id="' + id + '">Now loading...</a>';
	}, function() {
		alert(1);
	});
	this.qset.append(elm1);
	this.qset.append(elm2);
	this.qset.start(2000);
	this.last = 1;
}

Board.prototype.term = function() {
	this.qset.end();
}

Board.prototype.setDisable = function(f) {
	$('#content').prop('disabled', f);
	$('#btn_wri').prop('disabled', f);
}


Board.prototype.write = function (txt) {
	var _this = this;
	rpccall('Alterorg.WriteToBoard', [[this.address, txt, (this.last+1).toString()]], function(res) {
		$('#content').val('');
		_this.draw();
	},
	function(res, stat, err) {
		_this.draw();
		alert(err.message);
	},
	function() {
		alert('failed');
	});
}


Board.prototype.draw = function() {
	var _this = this;
	var fm = ''
	 + '<h1>Board Test</h1>'
	 + '<div class="form-group">'
	 + '<label>list</label>'
	 + '<div id="list">'
	 + '</div>'
	 + '</div>'
	 + '<div class="form-group">'
	 + '<label>Content</label>'
	 + '<textarea id="content" class="form-control" rows=10 side="30" disabled></textarea>'
	 + '<ul class="list-inline">'
	 + '<li><button class="btn btn-primary" id="btn_wri" disabled>Write</button></li>'
	 + '</ul>'
	 + '</div>'
	 + '</ul>';
	$('#main').html(fm);
	$('#btn_wri').click(function() {
		_this.write($('#content').val());
	});
	$M.breadCrumb([{name:'Home', link:'$M.home()'}, {name:'Assembly detail'}, {name:'Board'}]);
}
