import { Component, Input, signal, Signal, WritableSignal } from '@angular/core';
import { Article } from '../../domain/Article';

@Component({
  selector: 'klcs-cashdesk',
  imports: [],
  templateUrl: './cashdesk.component.html',
  styleUrl: './cashdesk.component.css'
})
export class CashdeskComponent {
  @Input() articles: Signal<Map<string, Article[]>> = signal(new Map());
}
