/*!
 * 对jquery.validate简单封装，针对不同的前端框架，默认参数值可能不一样
 * Created by lihaitao on 2017-9-27.
 */
(function ($) {
    $.fn.extend({
        sdvalidate: function (options) {
            // if nothing is selected, return nothing; can't chain anyway
            if (!this.length) {
                if (options && options.debug && window.console) {
                    console.warn("Nothing selected, can't validate, returning nothing.");
                }
                return;
            }
            var defaults={
                errorElement: 'span', //default input error message container
                errorClass: 'help-block help-block-error', // default input error message class
                focusInvalid: false, // do not focus the last invalid input
                ignore: "", // validate all fields including form hidden input
                rules: {},
                messages: {},
                errorPlacement: function (error, element) { // render error placement for each input type
                    console.log(element.parent(".input-group"))
                    
                    if (element.parent(".form-group").size() > 0) {
                        error.insertAfter(element.parent(".input-group"));
                    } else {
                        error.insertAfter(element); // for other inputs, just perform default behavior
                    }
                },
                invalidHandler: function (event, validator) { //display error alert on form submit
                    //验证不通过时
                },
                highlight: function (element) { // hightlight error inputs
                    $(element).closest('.form-group').removeClass('has-success').addClass('has-error');
                },
                unhighlight: function (element) {
                    $(element).closest('.form-group').removeClass('has-error');
                },
                success: function (label) {
                    label.closest('.form-group').removeClass('has-error').addClass("has-success");
                },
                submitHandler: null
            };
            var destOptions = $.extend({},defaults, options);
            //调用jquuery.validate插件方法
            this.validate(destOptions);
        },
    });
})(jQuery);