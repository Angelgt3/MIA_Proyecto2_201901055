import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  constructor(private router:Router,private http: HttpClient){}
  user: string = "";
  pass: string = "";
  id: string = "";
  Resultados=[  
    {
      nombre: "",
      result: "Resultado"  
    }
  ]
  compilar(){
    
    var comando = "login >user="+this.user+" >pwd="+this.pass+" >id="+this.id
    //this.http.post('http://34.224.222.47/ejecutar',comando).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/ejecutar',comando).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    this.user=""
    this.pass=""
    this.id=""
    console.log(this.Resultados[0].result)
      this.router.navigate(['inicio']);  
  }
  
  ir_inicio(){
    //Se deslogea tanbien
    var comando = "logout"
    //this.http.post('http://34.224.222.47/ejecutar',comando).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/ejecutar',comando).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    this.router.navigate(['inicio']);
  }
}
