#include<stddef.h>
#include<stdlib.h>
#include<stdio.h>
#include<assert.h>

typedef int content;

typedef struct node {
    content content;
    struct node* next;
    struct node* prev;
} node;

void append_to(node* parent, content val) {
    node* child = (node*) malloc(sizeof(node));
    child->content = val;
    child->prev = parent;

    if (parent->next == NULL){
        parent->next = child;
    //insert child behind parent append what was previously appended to parent
    } else { 
        node* parent_next = parent->next;
        parent->next = child;
        child->next = parent_next;
        parent_next->prev = child;
    }


}

int main() {
    node* head = NULL;
    head = (node*) malloc(sizeof(node));

    if (head == NULL) {
        printf("Head still NULL");
        return 1;
    }

    head->content = 2;
    assert(head->content == 2);

    append_to(head, 4);
    assert(head->next->content == 4);

    append_to(head, 3);
    assert(head->next->content == 3);
    assert(head->next->next->content == 4);


    return 0;
}
