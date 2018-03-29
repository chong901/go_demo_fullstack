let socket = io.connect();
let partOutput = [];
let jobOutput = [];

$(function(){
  $.ajax({
    url: '/api/v1/jobM/partQtyList'
  })
  .done((r) => {
    if(r) {
      $('#progressSummary').html(createProgressSummary(r));
    }
  })
  .fail(alertErr);

  $.ajax({
    url: '/api/v1/machine/summary'
  })
  .done((r) => {
    if(r.data){
      $('#machineTiles').html(createMachineTiles(r));
    }
  })
  .fail(alertErr);

  $('body').popover({ selector: '[data-toggle="popover"][data-trigger="hover"]', trigger: "hover" });
});

function createProgressSummary(res) {
  if(res.length == 0) {
    return '<div class="row"><div class="col-sm-3">' +
      '<div class="tile grey">' +
      `<h4 class="title"><a href="/jobM">${$('#strAddJob').text()}</a></h4>` +
      '</div>' +
      '</div></div>';
  } else {
    let html = '';
    let partWip = res;
    html = partWip.reduce((p, c) => {
      let wPercent = calcBarPercent(c.Output, c.Planned);
      partOutput.push({part: c.Id, total: c.Output, plan: c.Planned});
      return p + `<div class="col-sm-3">` +
      `<div class="tile grey" data-partWip="${c.Id}">` +
      `<h4 class="title">${c.Id}</h4>` +
      '<div class="progress">' +
      `<div class="progress-bar progress-bar-info" role="progressbar" aria-valuenow="${wPercent}" aria-valuemin="0" aria-value-max="100" style="min-width: 2em; width: ${wPercent}%;" data-skuProgress="${c.Id}">${wPercent}%</div>` +
      '</div>' +
      `<span class=jInfo-right><span data-skuInfo="${c.Id}">${c.Output}</span> / ${c.Planned}</span>` +
      '</div></div>';
    }, '');
    return '<div class="row">' + html + '</div>';
  }
}

function calcBarPercent(out, total) {
  if(out == 0) {
    return "0";
  } else if(out >= total || total == 0) {
    return "100";
  } else {
    return Math.floor((out/total)*100).toString();
  }
}

function pad(num) {
  return ("0"+num).slice(-2);
}

function hhmmss(secs) {
  let minutes = Math.floor(secs/60);
  secs = secs%60;
  let hours = Math.floor(minutes/60);
  minutes = minutes%60;
  return pad(hours)+":"+pad(minutes)+":"+pad(secs);
}

function getEstimate(ct, actual, plan) {
  let result = '--:--:--';
  let diff = plan - actual;
  let estSeconds = 0.0;
  if (diff > 0) {
    estSeconds = ct*diff;
    result = hhmmss(Math.round(estSeconds));
  }
  return result;
}

function createMachineTiles(res) {
  if(res.data.length == 0) {
    return '<div class="row"><div class="col-sm-3">' +
      '<div class="tile grey">' +
      `<h4 class="title"><a href="/machineM">${$('#strAddMachine').text()}</a></h4>` +
      '</div>' +
      '</div></div>';
  } else {
    let html = '';
    let stations = res.data;
    html = stations.reduce((p, c) => {
      let part = c.Job? c.Sku : $('#strNoJob').text();
      let output = c.Job? c.Output : '0';
      let quantity = c.Job? c.Qty : '0';
      let tileDiv = '';
      let statusTxt = '';
      let estimate = getEstimate(c.Ct, c.Output, c.Qty);
      if(c.Job) jobOutput.push({job: c.Job, total: output});
      switch(c.Status) {
        case 1:
          tileDiv = `<div class="tile green" data-machineName="${c.Name}">`;
          statusTxt = 'running';
          break;
        case 2:
          tileDiv = `<div class="tile red" data-machineName="${c.Name}">`;
          statusTxt = 'stopped';
          break;
        case 3:
          tileDiv = `<div class="tile orange" data-machineName="${c.Name}">`;
          statusTxt = 'line change';
          break;
        default:
          tileDiv = `<div class="tile blue" data-machineName="${c.Name}">`;
          statusTxt = 'idle';
      }
      let infoLink = c.Job?
        `<span class="mInfo-left"><a tabindex="0" class="info-a" data-toggle="popover" data-trigger="hover" data-container="body" data-html="true" title="${$('#strJob').text()} #${c.Job}" data-content="${$('#strStatus').text()}: ${statusTxt}<br>${$('#strRemaining').text()}: ${estimate}" data-placement="bottom">${part}</a></span>` :
        `<span class="mInfo-left">${part}</span>`;
      return p + '<div class="col-sm-3">' + tileDiv +
        `<h4 class="title">${c.Name}</h4>` + infoLink +
        `<span class="mInfo-right"><span data-machineJob="${c.Job}">${output}</span> / ${quantity}</span>` +
        '</div></div>';
    }, '');
    return '<div class="row">' + html + '</div>';
  }
}

socket.on('jobUpdate', function(data){
  //let msg = 'masterId: ' + data.masterId + ' sku: ' + data.sku + ' incVal: ' + data.incVal;
  //alert(msg);
  for(var i=0; i<partOutput.length; i++) {
    if(partOutput[i].part == data.sku) {
      partOutput[i].total += data.incVal;
      $('*[data-skuInfo="'+data.sku+'"]').html(partOutput[i].total);
      let progress = calcBarPercent(partOutput[i].total, partOutput[i].plan);
      $('*[data-skuProgress="'+data.sku+'"]').attr({
        "aria-valuenow": progress,
        "style": "min-width: 2em; width: "+progress+"%;"
      });
      $('*[data-skuProgress="'+data.sku+'"]').html(progress+'%');
      break;
    }
  }
  for(var j=0; j<jobOutput.length; j++) {
    if(jobOutput[j].job == data.masterId) {
      jobOutput[j].total += data.incVal;
      $('*[data-machineJob="'+data.masterId+'"]').html(jobOutput[j].total);
      break;
    }
  }
});