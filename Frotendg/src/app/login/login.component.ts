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
      nombre: "Archivo 1",
      result: "Resultado 1"  
    }
  ]
  compilar(){
    console.log(this.user)
    console.log(this.pass)
    console.log(this.id)
    var comando = "login >user="+this.user+" >pwd="+this.pass+" >id="+this.id
    this.http.post('http://localhost:3000/login',{nombre:"login",texto:comando}).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    this.user=""
    this.pass=""
    this.id=""
  }
  
  ir_inicio(){
    this.router.navigate(['inicio']);
  }
}
