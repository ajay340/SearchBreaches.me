const search = String(window.location.pathname).replace("/search=", "").replace(/%20/g, ' ');
const matchList = document.getElementById('row');
let breaches;


   const getBreachSearch = async (searchText) =>{
    let matches = [];
    if (searchText.length === 0) {
        matches = [];
        matchList.innerHTML = '';
    } else{
    url = '/searchq=' + searchText
    const res = await fetch(url);
    breaches = await res.json();
    matches = breaches.filter(breaches => {
        const regex = new RegExp(`^${searchText}`, 'gi');
        return breaches.Name_of_Covered_Entity || breaches.Industry || breaches.Type_of_Breach; 
       });
      
       // Clear when input or matches are empty
       if (searchText.length === 0) {
        matches = [];
        matchList.innerHTML = '';
       }
    }
      console.log(matches)
       outputHtml(matches);
      };

const outputHtml = matches => {
 if (matches.length > 0) {
  const html = matches
   .map(
    match => `
    <a href="/product/${match.ID}">
    <p><strong>${match.Name_of_Covered_Entity}
  </a></strong><br>${match.Summary}</p>
  <hr width="100%" size="1" align="left">
  <br>
    `
   )
   .join('');
  matchList.innerHTML = html;
 }
};
window.addEventListener('DOMContentLoaded', () => getBreachSearch(search));
