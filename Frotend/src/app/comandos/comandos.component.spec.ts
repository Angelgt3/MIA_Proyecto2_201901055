import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ComandosComponent } from './comandos.component';

describe('ComandosComponent', () => {
  let component: ComandosComponent;
  let fixture: ComponentFixture<ComandosComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ComandosComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ComandosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
