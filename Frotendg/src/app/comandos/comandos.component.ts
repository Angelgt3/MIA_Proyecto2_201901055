import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-comandos',
  templateUrl: './comandos.component.html',
  styleUrls: ['./comandos.component.css']
})
export class ComandosComponent {
  constructor(private router:Router,private http: HttpClient){}
  comand: string = "";
  
  Resultados=[  
    {
      nombre: "Archivo 1",
      result: "Resultado 1",
      reporte: "/sb.jpg"
    }
  ]

  compilar(){
    console.log(this.comand);
    //this.http.post('http://34.224.222.47/ejecutar',comando).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/reporte_disk',this.comand).subscribe((response:any)=>{
      this.Resultados[0].reporte=response
    });
    this.comand="";
  }
  salir(){
    //Se deslogea tanbien
    var comando = "logout"
    //this.http.post('http://34.224.222.47/ejecutar',comando).subscribe((response:any)=>{
    this.http.post('http://localhost:3000/ejecutar',comando).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    this.router.navigate(['login']);
  }
}
