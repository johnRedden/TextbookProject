    $.wait = function( callback, seconds){
       return window.setTimeout( callback, seconds * 1000 );
    }

// no images ckeditor config
var noImagesConfig = {
  	extraPlugins: 'mathjax',
  	mathJaxLib: 'https://cdn.mathjax.org/mathjax/2.6-latest/MathJax.js?config=TeX-AMS_HTML',
   	removePlugins: 'forms'
 };

 // images ckeditor requires an editity id
 function stdImageConfig(argID){
    var stdConfig = {
            extraPlugins: 'mathjax',
            mathJaxLib: 'https://cdn.mathjax.org/mathjax/2.6-latest/MathJax.js?config=TeX-AMS_HTML',
            filebrowserImageBrowseUrl: '/image/browser?action=browseImageURL&oid='+argID,
            filebrowserImageUploadUrl: '/api/ckeditor/create?action=uploadImageURL&oid='+argID,
            filebrowserBrowseUrl: '/image/browser?action=browseURL&oid='+argID,
            filebrowserUploadUrl: '/api/ckeditor/create?action=uploadURL&oid='+argID,
            removePlugins: 'forms',
            skin: 'icy_orange,/public/icy_orange/'
    };
    //add any special html chars here
    stdConfig.specialChars =[ "α","β","ψ","δ","ε","φ","γ","η","ι","ξ","κ","λ","μ","ν","ο","π",";","ρ","σ","τ","θ","ω","ς","χ","υ","ζ",'°','ℝ','ℂ','ℕ','ℚ','ℤ','∀','∁','∂','∃','∄','∅','Δ','∇','∈','∉','ε','∋','∌','∍','∎','Π','∐','∠','∡','Σ','−','±','∓','∔','∕','','∗','∘','∙','√','∛','∜','∝','∞','∟','∠','∡','∢','∣','∤','∥','∦','∧','∨','∩','∪','∫','∬','∭','∮','∯','∰','∱','∲','∳','∴','∵','∶','∷','∸','∹','∺','∻','∼','∽','∾','∿','≀','≁','≂','≃','≄','≅','≆','≇','≈','≉','≊','≋','≌','≍','≎','≏','≐','≑','≒','≓','≔','≕','≖','≗','≘','≙','≚','≛','≜','≝','≞','≟','≠','≡','≢','≣','≤','≥','≦','≧','≨','≩','≪','≫','≬','≭','≮','≯','≰','≱','≲','≳','≴','≵','≶','≷','','','≺','≻','≼','≽','≾','≿','⊀','⊁','⊂','⊃','⊄','⊅','⊆','⊇','⊈','⊉','⊊','⊋','⊌','⊍','⊎','⊏','⊐','⊑','⊓','⊔','⊕','⊖','⊗','⊘','⊙','⊚','⊛','⊜','⊝','⊞','⊟','⊠','⊡','⊢','⊣','⊤','⊥','⊦','⊧','⊨','⊩','⊪','⊫','⊬','⊭','⊮','⊯','⊰','⊱','⊲','⊳','⊴','⊵','','','','⊹','⊺','⊻','⊼','⊽','⊿','⋀','⋁','⋂','⋃','⋄','⋅','⋆','⋇','','','⋍','⋎','⋏','⋐','⋑','⋒','⋓','⋔','⋕','⋖','⋗','⋘','⋙','⋚','⋛','⋜','⋝','⋞','⋟','⋠','⋡','⋢','⋣','⋤','⋥','⋦','⋧','⋨','⋩','⋪','⋫','⋬','⋭','⋮','⋯','⋰','⋱',];
    return stdConfig;
 }

// code stolen to overried the save event
CKEDITOR.plugins.registered['save'] = {
  init: function (editor) {
     var command = editor.addCommand('save',
     {
          modes: { wysiwyg: 1, source: 1 },
          exec: function (editor) { // Add here custom function for the save button
            //alert('You clicked the save button in CKEditor toolbar!');
            $('.saveBtn').trigger('click');
          }
     });
     editor.ui.addButton('Save', { label: 'Save', command: 'save' });
  }
}

CKEDITOR.config.contentsCss = '/public/css/myeditor.css';

function ensBreadcrumb(entity,id){
  //URL is /api/parent/:Kind/:ID
  var promise = $.get('/api/parent/'+entity+'/'+id,function(data, status){
      var x = $.parseJSON(data).Results;
     console.log(x);
  });
  return promise;

}