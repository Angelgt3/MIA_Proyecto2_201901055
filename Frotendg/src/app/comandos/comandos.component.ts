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
      result: "Resultado 1"  
    }
  ]

  compilar(){
    console.log(this.comand);
    
    this.http.post('http://localhost:3000/ejecutar_comando',{contenido:this.comand}).subscribe((response:any)=>{
      this.Resultados[0]=response
    });
    this.comand="";
  }
  salir(){
    this.router.navigate(['login']);
  }
}
