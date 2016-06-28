

$(document).on('ready',function(){

//menuBar   **********************************************
    $('#fontChange').click(function(){
        $('#submenu1').toggleClass('hidden');
        $('#submenu2').addClass('hidden');
    })
    $('#increaseFont').click(function(){
        var x = parseInt($('.bookPage').css('font-size'));
        if(x<48){
            $('.bookPage').css('font-size', x+2+'px');
        }
        $(this).blur();
    });
    $('#decreaseFont').click(function(){
        var x = parseInt($('.bookPage').css('font-size'));
        if(x>8)
            $('.bookPage').css('font-size', x-2+'px');
        $(this).blur();
    });
    $('#imageChange').click(function(){
        $('#submenu2').toggleClass('hidden');
        $('#submenu1').addClass('hidden');
    })            
    $('#increaseImages').click(function(){
        var imgs = $('img');
        $.each(imgs,function(index, arg){
            var w = $(arg).width();
            w+=20;
            $(arg).width(w+'px');
        })
        console.log(imgs);
        $(this).blur();
    });

    $('#decreaseImages').click(function(){
        var imgs = $('img');
        $.each(imgs,function(index, arg){
            var w = $(arg).width();
            w-=20;
            $(arg).width(w+'px');
        })
        $(this).blur();
    });
    $('#print').click(function(){
        window.print();
        $(this).blur();
    });
// end menubar  ****************************************

    
});



