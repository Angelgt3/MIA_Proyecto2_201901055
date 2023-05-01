import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-inicio',
  templateUrl: './inicio.component.html',
  styleUrls: ['./inicio.component.css']
  },
)
export class InicioComponent {
  
  constructor(private router:Router,private http: HttpClient){}
  
  path: string = "";
  code: string = "";
  Resultados=[  
    {
      nombre: "",
      result: "Resultado"  
    }
  ]
  compilar(){
    //console.log(this.code);
    
    //this.http.post('http://34.224.222.47/ejecutar',this.code).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/ejecutar',this.code).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    console.log(this.Resultados[0].result);
    //this.path="";
    //this.code="";
  }
  
  ir_login(){
    this.router.navigate(['login']);
  }
  salir(){
    var comando = "logout"
    //this.http.post('http://34.224.222.47/ejecutar',comando).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/ejecutar',comando).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
  }
}