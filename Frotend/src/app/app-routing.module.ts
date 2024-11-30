import { Component, NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ComandosComponent } from './comandos/comandos.component';
import { InicioComponent } from './inicio/inicio.component';
import { LoginComponent } from './login/login.component';

const routes: Routes = [
  {
    path:'inicio',
    component: InicioComponent
  },
  {
    path:'login',
    component: LoginComponent
  },
  {
    path:'comandos',
    component: ComandosComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
