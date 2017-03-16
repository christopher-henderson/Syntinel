function getQueryVariable(variable) {
   var query = window.location.search.substring(1);
   var vars = query.split("&");
   for (var i=0;i<vars.length;i++) {
           var pair = vars[i].split("=");
           if(pair[0] == variable){return pair[1];}
   }
}

function escapeNewLineChars(valueToEscape) {
   if (valueToEscape != null && valueToEscape != "") {
      return valueToEscape.replace(/\n/g, "\\n");
   } else {
      return valueToEscape;
   } 
}